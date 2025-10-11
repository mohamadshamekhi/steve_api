package mongodbcomment

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/abolfazlnorzad/steve/entity"
	commentparam "github.com/abolfazlnorzad/steve/param/comment"
	commonparam "github.com/abolfazlnorzad/steve/param/common"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	commentservice "github.com/abolfazlnorzad/steve/service/comment"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d DB) List(ctx context.Context, f commonparam.FilterRequest, p commonparam.PaginationRequest) (commentservice.ListResult, error) {
	const op = "mongodbcomment.List"

	filter := bson.M{
		"deleted_at": bson.M{"$exists": false},
	}
	for _, filterItem := range f.Filter {
		for key, value := range filterItem {

			valueStr := fmt.Sprintf("%v", value)

			if valueStr == "" {
				continue
			}

			switch key {

			case "is_approved":
				if boolValue, err := strconv.ParseBool(valueStr); err == nil {
					filter[key] = boolValue
				}

			default:
				filter[key] = value
			}
		}
	}

	offset := p.GetOffset()
	limit := p.GetPageSize()
	findOptions := options.Find().
		SetSkip(int64(offset)).
		SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	fmt.Println("filter", filter)
	totalCount, err := d.collection.CountDocuments(ctx, filter)
	if err != nil {
		return commentservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	cursor, err := d.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return commentservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	defer cursor.Close(ctx)

	var comments []entity.Comment
	if err := cursor.All(ctx, &comments); err != nil {
		return commentservice.ListResult{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return commentservice.ListResult{
		Comments:   comments,
		TotalCount: uint(totalCount),
	}, nil
}

func (d DB) GetApprovedCommentsWithChildren(ctx context.Context, commentableType string, commentableID string, pagination commonparam.PaginationRequest) (commentservice.StorefrontListResult, error) {
	// 1. شمارش کل ریشه‌ها (برای Pagination)
	countPipeline := mongo.Pipeline{
		{
			{"$match", bson.M{
				"commentable_id":   commentableID,
				"commentable_type": commentableType,
				"parent_id":        bson.M{"$exists": false},
				"is_approved":      true,
				"deleted_at":       bson.M{"$exists": false},
			}},
		},
		{{"$count", "totalCount"}},
	}

	countCursor, err := d.collection.Aggregate(ctx, countPipeline)
	if err != nil {
		return commentservice.StorefrontListResult{}, fmt.Errorf("failed to count roots: %w", err)
	}
	defer countCursor.Close(ctx)

	var countResult []bson.M
	countCursor.All(ctx, &countResult)

	totalCount := uint(0)
	if len(countResult) > 0 {
		if count, ok := countResult[0]["totalCount"].(int32); ok {
			totalCount = uint(count)
		}
	}

	// اگر کامنت ریشه‌ای وجود ندارد، با TotalCount صفر برگردان
	if totalCount == 0 {
		return commentservice.StorefrontListResult{TotalCount: 0}, nil
	}

	// 2. واکشی IDهای ریشه‌های صفحه‌بندی‌شده (کوئری اول)
	pipelineRootIDs := mongo.Pipeline{
		{
			{"$match", bson.M{
				"commentable_id":   commentableID,
				"commentable_type": commentableType,
				"parent_id":        bson.M{"$exists": false},
				"is_approved":      true,
				"deleted_at":       bson.M{"$exists": false},
			}},
		},
		{{"$sort", bson.D{{Key: "created_at", Value: -1}}}},
		{{"$skip", pagination.GetOffset()}},
		{{"$limit", pagination.GetPageSize()}},
		{{"$project", bson.M{"_id": 1}}}, // فقط IDها را برگردان
	}

	rootIDCursor, err := d.collection.Aggregate(ctx, pipelineRootIDs)
	if err != nil {
		return commentservice.StorefrontListResult{}, fmt.Errorf("failed to fetch root IDs: %w", err)
	}
	defer rootIDCursor.Close(ctx)

	var rootDocs []bson.M
	rootIDCursor.All(ctx, &rootDocs)

	var rootIDs []string
	for _, doc := range rootDocs {
		if id, ok := doc["_id"].(string); ok {
			rootIDs = append(rootIDs, id)
		}
	}

	orConditions := bson.A{
		bson.M{"_id": bson.M{"$in": rootIDs}},           // <--- رفع خطای ۱: تعریف صریح
		bson.M{"ancestors._id": bson.M{"$in": rootIDs}}, // <--- رفع خطای ۲: تعریف صریح
	}

	filterAllComments := bson.M{
		"$or":              orConditions,
		"commentable_id":   commentableID,
		"commentable_type": commentableType,
		"is_approved":      true,
		"deleted_at":       bson.M{"$exists": false},
	}

	var flatComments []entity.Comment // فرض می‌کنیم entity.Comment ساختار اصلی بدون فیلد Children است
	cursor, err := d.collection.Find(ctx, filterAllComments)
	if err != nil {
		return commentservice.StorefrontListResult{}, fmt.Errorf("failed to fetch flat comments: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &flatComments); err != nil {
		return commentservice.StorefrontListResult{}, fmt.Errorf("failed to decode flat comments: %w", err)
	}

	// 4. ساختاردهی درختی و بازگرداندن نتایج (انتقال به سرویس)
	// در عمل، بهتر است این مرحله در لایه سرویس انجام شود، اما برای سادگی، در اینجا آن را فراخوانی می‌کنیم.
	paginatedRoots := BuildCommentTree(flatComments, rootIDs)

	return commentservice.StorefrontListResult{
		Comments:   paginatedRoots,
		TotalCount: totalCount,
	}, nil
}

// service/comment_utils.go یا service/comment_service.go

// فرض: CommentWithChildren ساختار نهایی شماست که Children آن []*CommentWithChildren است.

func BuildCommentTree(flatComments []entity.Comment, paginatedRootIDs []string) []commentparam.CommentWithChildren {
	commentMap := make(map[string]*commentparam.CommentWithChildren)

	// 1. نگاشت: تبدیل به ساختار خروجی و پر کردن نقشه
	for _, c := range flatComments {
		// ... تبدیل entity.Comment به *commentparam.CommentWithChildren ...
		// استفاده از strings.TrimSpace برای ParentID الزامی است.

		comment := &commentparam.CommentWithChildren{
			ID:              c.ID,
			Body:            c.Body,
			CommentableID:   c.CommentableID,
			CommentableType: c.CommentableType,
			CommenterID:     c.CommenterID,
			CommenterType:   c.CommenterType,
			CommenterName:   c.CommenterName,
			ParentID:        strings.TrimSpace(c.ParentID), // تمیز کردن ParentID
			Ancestors:       c.Ancestors,
			Level:           c.Level,
			IsApproved:      c.IsApproved,
			CreatedAt:       c.CreatedAt,
			UpdatedAt:       c.UpdatedAt,
			Children:        make([]*commentparam.CommentWithChildren, 0),
		}
		commentMap[c.ID] = comment
	}

	// 2. اتصال: برقراری پیوند والد و فرزند
	for _, c := range commentMap {
		if c.ParentID != "" {
			if parent, ok := commentMap[c.ParentID]; ok {
				// اگر از Children []CommentWithChildren استفاده می‌کنید، نیاز به تبدیل به مقدار است
				// اگر از Children []*CommentWithChildren استفاده می‌کنید، باید c را اضافه کنید
				parent.Children = append(parent.Children, c) // فرض بر این است که Children شما مقدار است
			}
		}
	}

	// 3. فیلتر کردن نهایی: برگرداندن فقط ریشه‌هایی که ID آن‌ها در paginatedRootIDs است
	finalRoots := make([]commentparam.CommentWithChildren, 0, len(paginatedRootIDs))
	rootIDSet := make(map[string]bool)
	for _, id := range paginatedRootIDs {
		rootIDSet[id] = true
	}

	for _, c := range commentMap {
		if _, isRoot := rootIDSet[c.ID]; isRoot {
			// از آنجایی که ما فقط به ریشه‌های صفحه‌بندی‌شده نیاز داریم،
			// و در مرحله اتصال، تمام فرزندان به آن‌ها متصل شده‌اند،
			// می‌توانیم مستقیم آن‌ها را به لیست نهایی اضافه کنیم.
			finalRoots = append(finalRoots, *c)
		}
	}

	return finalRoots
}
