package articleservice

import (
	"context"
	"errors"
	"fmt"
	"github.com/abolfazlnorzad/steve/entity"
	articleparam "github.com/abolfazlnorzad/steve/param/article"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/abolfazlnorzad/steve/pkg/slug"
)

type BackofficeService struct {
	repo     ArticleReadWriter
	mediaSvc MediaService
}

func NewBackofficeService(repo ArticleReadWriter, mediaSvc MediaService) BackofficeService {
	return BackofficeService{repo: repo, mediaSvc: mediaSvc}
}

func (s BackofficeService) Create(ctx context.Context, req articleparam.AddArticleRequest) (articleparam.AddArticleResponse, error) {
	const op = "articleservice.Create"

	a := entity.Article{
		Title:       req.Title,
		Body:        req.Body,
		IsPublished: req.IsPublished,
		MediaID:     req.MediaID,
	}
	slug := slug.GenerateSlug(a.Title)
	counter := 1
	for {
		existingArticle, err := s.repo.FindBySlug(ctx, slug)
		if err != nil {
			var rErr *richerror.RichError
			if errors.As(err, &rErr) && rErr.Kind() == richerror.KindNotFound {
				break
			}
			return articleparam.AddArticleResponse{}, richerror.New(op).WithWarpError(err)
		}

		if existingArticle.ID == "" {
			break
		}

		slug = fmt.Sprintf("%s-%d", slug, counter)
		counter++
	}

	a.Slug = slug
	a.BeforeCreate()

	if err := s.repo.Create(ctx, a); err != nil {
		return articleparam.AddArticleResponse{}, richerror.New(op).WithWarpError(err)
	}
	res, err := s.mediaSvc.GetURLByMediaID(ctx, mediaparam.GetUrlByIDRequest{ID: a.MediaID})
	if err != nil {
		return articleparam.AddArticleResponse{}, richerror.New(op).WithWarpError(err)
	}
	articleInfo := articleparam.ArticleInfo{
		ID:          a.ID,
		Title:       a.Title,
		Slug:        a.Slug,
		Body:        a.Body,
		IsPublished: a.IsPublished,
		MediaUrl:    res.URL,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
	return articleparam.AddArticleResponse{
		Article: articleInfo,
	}, nil

}

func (s BackofficeService) Update(ctx context.Context, req articleparam.UpdateArticleRequest) (articleparam.UpdateArticleResponse, error) {
	const op = "articleservice.Update"
	article, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return articleparam.UpdateArticleResponse{}, richerror.New(op).WithWarpError(err)
	}

	var oldMediaIDToDelete string

	if req.Title != nil && *req.Title != article.Title {
		article.Title = *req.Title

		slug := slug.GenerateSlug(article.Title)
		counter := 1
		for {
			existingArticle, err := s.repo.FindBySlug(ctx, slug)
			if err != nil {
				var rErr *richerror.RichError
				if errors.As(err, &rErr) && rErr.Kind() == richerror.KindNotFound {
					break
				}
				return articleparam.UpdateArticleResponse{}, richerror.New(op).WithWarpError(err)
			}

			if existingArticle.ID == "" {
				break
			}

			slug = fmt.Sprintf("%s-%d", slug, counter)
			counter++
		}

		article.Slug = slug
	}

	if req.Body != nil {
		article.Body = *req.Body
	}
	if req.IsPublished != nil {
		article.IsPublished = *req.IsPublished
	}

	if req.MediaID != nil && *req.MediaID != article.MediaID {

		_, mediaErr := s.mediaSvc.GetURLByMediaID(ctx, mediaparam.GetUrlByIDRequest{ID: *req.MediaID})
		if mediaErr != nil {
			var rErr *richerror.RichError
			if errors.As(mediaErr, &rErr) && rErr.Kind() == richerror.KindNotFound {
				return articleparam.UpdateArticleResponse{}, richerror.New(op).
					WithKind(richerror.KindInvalid).
					WithMessageKey("media.not_found").
					WithMeta(map[string]any{"media_id": *req.MediaID})
			}
			return articleparam.UpdateArticleResponse{}, richerror.New(op).WithWarpError(mediaErr)
		}

		oldMediaIDToDelete = article.MediaID

		article.MediaID = *req.MediaID
	}

	if uErr := s.repo.Update(ctx, article); uErr != nil {
		return articleparam.UpdateArticleResponse{}, richerror.New(op).WithWarpError(uErr).WithKind(richerror.KindUnexpected)
	}
	if oldMediaIDToDelete != "" {
		go func() {
			r, err := s.mediaSvc.FindMediaByID(context.Background(), mediaparam.GetMediaByIDRequest{ID: oldMediaIDToDelete})
			if err != nil {
				logger.L().Error("error in find old media", oldMediaIDToDelete)
			}
			if _, err := s.mediaSvc.DeleteMedia(context.Background(), mediaparam.DeleteMediaRequest{Media: r.Media}); err != nil {
				logger.L().Error("ERROR: background media deletion failed for ID %s: %v\n", oldMediaIDToDelete, err)
			}
		}()
	}

	res, err := s.mediaSvc.GetURLByMediaID(ctx, mediaparam.GetUrlByIDRequest{ID: article.MediaID})
	if err != nil {
		return articleparam.UpdateArticleResponse{}, richerror.New(op).WithWarpError(err)
	}
	articleInfo := articleparam.ArticleInfo{
		ID:          article.ID,
		Title:       article.Title,
		Slug:        article.Slug,
		Body:        article.Body,
		IsPublished: article.IsPublished,
		MediaUrl:    res.URL,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}

	return articleparam.UpdateArticleResponse{
		Article: articleInfo,
	}, nil
}

func (s BackofficeService) Delete(ctx context.Context, req articleparam.DeleteArticleRequest) (articleparam.DeleteArticleResponse, error) {
	const op = "articleservice.Delete"
	article, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return articleparam.DeleteArticleResponse{}, richerror.New(op).WithWarpError(err)
	}
	if !article.DeletedAt.IsZero() {
		return articleparam.DeleteArticleResponse{}, richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
	}
	if dErr := s.repo.Delete(ctx, article); dErr != nil {
		return articleparam.DeleteArticleResponse{}, richerror.New(op).WithWarpError(dErr)
	}
	// todo-write background job for delete media after one day and children's media
	//go func() {
	//	if _, err := s.mediaSvc.DeleteMedia(context.Background(), mediaparam.DeleteMediaRequest{ID: brand.MediaID}); err != nil {
	//		logger.L().Error("ERROR: background media deletion failed for ID %s: %v\n", brand.MediaID, err)
	//	}
	//}()

	return articleparam.DeleteArticleResponse{}, nil
}

func (s BackofficeService) GetArticle(ctx context.Context, req articleparam.GetArticleRequest) (articleparam.GetArticleResponse, error) {
	const op = "articleservice.GetArticle"
	article, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return articleparam.GetArticleResponse{}, richerror.New(op).WithWarpError(err)
	}
	if !article.DeletedAt.IsZero() {
		return articleparam.GetArticleResponse{}, richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
	}
	res, err := s.mediaSvc.GetURLByMediaID(ctx, mediaparam.GetUrlByIDRequest{ID: article.MediaID})
	if err != nil {
		logger.L().Error(fmt.Sprintf("ERROR: article media get failed for ID %s\n", article.MediaID))
	}
	a := articleparam.ArticleInfo{
		ID:          article.ID,
		Title:       article.Title,
		Slug:        article.Slug,
		Body:        article.Body,
		IsPublished: article.IsPublished,
		MediaUrl:    res.URL,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}
	return articleparam.GetArticleResponse{
		Article: a,
	}, nil
}

func (s BackofficeService) GetAll(ctx context.Context, req articleparam.AllArticleRequest) (articleparam.AllArticleResponse, error) {
	const op = "articleservice.GetAll"
	resp, err := s.repo.All(ctx, req.Filter, req.Pagination)
	if err != nil {
		return articleparam.AllArticleResponse{}, richerror.New(op).WithWarpError(err)
	}
	var articlesInfo []articleparam.ArticleInfo

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
		articlesInfo = append(articlesInfo, articleparam.ArticleInfo{
			ID:          a.ID,
			Title:       a.Title,
			Slug:        a.Slug,
			Body:        a.Body,
			IsPublished: a.IsPublished,
			MediaUrl:    urlMap.URLs[a.MediaID],
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
		})
	}

	return articleparam.AllArticleResponse{
		Pagination: struct {
			PageSize   uint `json:"page_size"`
			PageNumber uint `json:"page_number"`
			Total      uint `json:"total"`
		}{PageSize: req.Pagination.GetPageSize(), PageNumber: req.Pagination.GetPageNumber(), Total: uint(resp.TotalCount)},
		Articles: articlesInfo,
	}, nil
}
