package articleservice

import (
	"context"
	"fmt"

	articleparam "github.com/abolfazlnorzad/steve/param/article"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	searchparam "github.com/abolfazlnorzad/steve/param/search"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type StorefrontService struct {
	aRepo       ArticleReader
	mediaSvc    MediaService
	bRepo       BackOfficeReader
	viewCounter ArticleViewCounter
	viewDedup   ArticleViewDeduplicator
}

func NewStorefrontService(
	articleRepo ArticleReader,
	mediaSvc MediaService,
	backOffice BackOfficeReader,
	viewCounter ArticleViewCounter,
	viewDedup ArticleViewDeduplicator,
) StorefrontService {
	return StorefrontService{
		aRepo:       articleRepo,
		mediaSvc:    mediaSvc,
		bRepo:       backOffice,
		viewCounter: viewCounter,
		viewDedup:   viewDedup,
	}
}

// registerView records a unique-per-visitor view for an article and returns
// the delta (0 or 1) that should be added to the count already read from
// storage, so the response reflects this view too. It never fails the
// request: view tracking is best-effort, so errors are only logged.
func (s StorefrontService) registerView(ctx context.Context, articleID, visitorID string) int64 {
	const op = "ArticleService.registerView"
	if articleID == "" || visitorID == "" {
		return 0
	}

	isNewViewer, err := s.viewDedup.AddViewer(ctx, articleID, visitorID)
	if err != nil {
		logger.L().Error(fmt.Sprintf("%s: dedup check failed for article %s", op, articleID), "error", err)
		return 0
	}
	if !isNewViewer {
		return 0
	}

	if err := s.viewCounter.IncrementViewCount(ctx, articleID); err != nil {
		logger.L().Error(fmt.Sprintf("%s: increment failed for article %s", op, articleID), "error", err)
		return 0
	}

	return 1
}

func (s StorefrontService) GetSingleStorefrontArticle(ctx context.Context, req articleparam.GetStorefrontArticleRequest) (articleparam.GetStorefrontArticleResponse, error) {
	const op = "ArticleService.GetSingleStorefrontArticle"
	article, err := s.aRepo.FindByID(ctx, req.ID)
	if err != nil {
		return articleparam.GetStorefrontArticleResponse{}, richerror.New(op).WithWarpError(err)
	}
	if !article.DeletedAt.IsZero() || !article.IsPublished {
		return articleparam.GetStorefrontArticleResponse{}, richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
	}

	backOfficeInfo, bErr := s.bRepo.GetBackOfficeUserByID(ctx, article.BackOfficeID)
	if bErr != nil {
		return articleparam.GetStorefrontArticleResponse{}, richerror.New(op).WithWarpError(bErr)
	}

	res, err := s.mediaSvc.GetURLByMediaID(ctx, mediaparam.GetUrlByIDRequest{ID: article.MediaID})
	if err != nil {
		logger.L().Error(fmt.Sprintf("ERROR: article media get failed for ID %s\n", article.MediaID))
	}
	a := articleparam.StorefrontArticleDetail{
		ID:              article.ID,
		Title:           article.Title,
		Slug:            article.Slug,
		Body:            article.Body,
		IsPublished:     article.IsPublished,
		Tags:            article.Tags,
		MetaDescription: article.MetaDescription,
		BackOfficeID:    article.BackOfficeID,
		BackOfficeName:  backOfficeInfo.Name,
		MediaUrl:        res.URL,
		Alt:             article.Alt,
		ViewCount:       article.ViewCount,
		CreatedAt:       article.CreatedAt,
		UpdatedAt:       article.UpdatedAt,
	}
	a.ViewCount += s.registerView(ctx, article.ID, req.VisitorID)
	return articleparam.GetStorefrontArticleResponse{
		Article: a,
	}, nil
}

func (s StorefrontService) GetSingleBySlugStorefrontArticle(ctx context.Context, req articleparam.GetBySlugStorefrontArticleRequest) (articleparam.GetBySlugStorefrontArticleResponse, error) {
	const op = "ArticleService.GetSingleBySlugStorefrontArticle"
	article, err := s.aRepo.FindBySlug(ctx, req.Slug)
	if err != nil {
		return articleparam.GetBySlugStorefrontArticleResponse{}, richerror.New(op).WithWarpError(err)
	}
	if !article.DeletedAt.IsZero() || !article.IsPublished {
		return articleparam.GetBySlugStorefrontArticleResponse{}, richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
	}

	backOfficeInfo, bErr := s.bRepo.GetBackOfficeUserByID(ctx, article.BackOfficeID)
	if bErr != nil {
		return articleparam.GetBySlugStorefrontArticleResponse{}, richerror.New(op).WithWarpError(bErr)
	}

	res, err := s.mediaSvc.GetURLByMediaID(ctx, mediaparam.GetUrlByIDRequest{ID: article.MediaID})
	if err != nil {
		logger.L().Error(fmt.Sprintf("ERROR: article media get failed for ID %s\n", article.MediaID))
	}
	a := articleparam.StorefrontArticleDetail{
		ID:              article.ID,
		Title:           article.Title,
		Slug:            article.Slug,
		Body:            article.Body,
		IsPublished:     article.IsPublished,
		Tags:            article.Tags,
		MetaDescription: article.MetaDescription,
		BackOfficeID:    article.BackOfficeID,
		BackOfficeName:  backOfficeInfo.Name,
		MediaUrl:        res.URL,
		Alt:             article.Alt,
		ViewCount:       article.ViewCount,
		CreatedAt:       article.CreatedAt,
		UpdatedAt:       article.UpdatedAt,
	}
	a.ViewCount += s.registerView(ctx, article.ID, req.VisitorID)
	return articleparam.GetBySlugStorefrontArticleResponse{
		Article: a,
	}, nil
}

func (s StorefrontService) GetAllStorefrontArticles(ctx context.Context, req articleparam.AllStorefrontArticleRequest) (articleparam.AllStorefrontArticleResponse, error) {
	const op = "ArticleService.GetAllStorefrontArticles"
	resp, err := s.aRepo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return articleparam.AllStorefrontArticleResponse{}, richerror.New(op).WithWarpError(err)
	}
	var articlesInfo []articleparam.StorefrontArticle

	mediaIDs := make([]string, 0, len(resp.Articles))
	mediaIDSet := make(map[string]struct{})
	for _, c := range resp.Articles {
		if c.MediaID != "" {
			if _, exists := mediaIDSet[c.MediaID]; !exists {
				mediaIDSet[c.MediaID] = struct{}{}
				mediaIDs = append(mediaIDs, c.MediaID)
			}
		}
	}

	urlMap, urlErr := s.mediaSvc.GetURLsByMediaIDs(ctx, mediaparam.GetUrlsByIDsRequest{IDs: mediaIDs})
	if urlErr != nil {
		logger.L().Error("warning: could not get media urls", "error", urlErr)
	}

	for _, a := range resp.Articles {

		//todo-یک متد جدید برای گرفتن همزمان چند انتیتی با استفاده از چند ای دی
		backOfficeInfo, bErr := s.bRepo.GetBackOfficeUserByID(ctx, a.BackOfficeID)
		if bErr != nil {
			return articleparam.AllStorefrontArticleResponse{}, richerror.New(op).WithWarpError(err)
		}

		articlesInfo = append(articlesInfo, articleparam.StorefrontArticle{
			ID:             a.ID,
			Title:          a.Title,
			Slug:           a.Slug,
			IsPublished:    a.IsPublished,
			Tags:           a.Tags,
			BackOfficeID:   a.BackOfficeID,
			BackOfficeName: backOfficeInfo.Name,
			MediaUrl:       urlMap.URLs[a.MediaID],
			Alt:            a.Alt,
			ViewCount:      a.ViewCount,
			CreatedAt:      a.CreatedAt,
			UpdatedAt:      a.UpdatedAt,
		})
	}

	return articleparam.AllStorefrontArticleResponse{
		Pagination: struct {
			PageSize   uint `json:"page_size"`
			PageNumber uint `json:"page_number"`
			Total      uint `json:"total"`
		}{PageSize: req.Pagination.GetPageSize(), PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)},
		Articles: articlesInfo,
	}, nil

}

func (s StorefrontService) GetSearchArticles(ctx context.Context, req articleparam.StorefrontGetSearchArticleRequest) (articleparam.StorefrontGetSearchArticleResponse, error) {
	const op = "ArticleService.GetSearchArticles"

	resp, err := s.aRepo.GetSearchArticles(ctx, req.Search)
	if err != nil {
		return articleparam.StorefrontGetSearchArticleResponse{}, richerror.New(op).WithWarpError(err)
	}

	var list []searchparam.StorefrontArticleSearchList

	mediaIDs := make([]string, 0, len(resp.Articles))
	mediaIDSet := make(map[string]struct{})
	for _, c := range resp.Articles {
		if c.MediaID != "" {
			if _, exists := mediaIDSet[c.MediaID]; !exists {
				mediaIDSet[c.MediaID] = struct{}{}
				mediaIDs = append(mediaIDs, c.MediaID)
			}
		}
	}

	urlMap, urlErr := s.mediaSvc.GetURLsByMediaIDs(ctx, mediaparam.GetUrlsByIDsRequest{IDs: mediaIDs})
	if urlErr != nil {
		logger.L().Error("warning: could not get media urls", "error", urlErr)
	}

	for _, article := range resp.Articles {
		list = append(list, searchparam.StorefrontArticleSearchList{
			ID:          article.ID,
			Title:       article.Title,
			Slug:        article.Slug,
			IsPublished: article.IsPublished,
			MediaUrl:    urlMap.URLs[article.MediaID],
			Alt:         article.Alt,
			CreatedAt:   article.CreatedAt,
			UpdatedAt:   article.UpdatedAt,
		})
	}

	return articleparam.StorefrontGetSearchArticleResponse{
		Articles: list,
	}, nil

}
