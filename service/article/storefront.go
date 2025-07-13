package articleservice

import (
	"context"
	articleparam "github.com/abolfazlnorzad/steve/param/article"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

type StorefrontService struct {
	repo     ArticleReader
	mediaSvc MediaService
}

func NewStorefrontService(repo ArticleReader, mediaSvc MediaService) StorefrontService {
	return StorefrontService{
		repo:     repo,
		mediaSvc: mediaSvc,
	}
}

func (s StorefrontService) GetSingleStorefrontArticle(ctx context.Context, req articleparam.GetStorefrontArticleRequest) (articleparam.GetStorefrontArticleResponse, error) {
	const op = "ArticleService.Get"
	article, err := s.repo.FindByID(ctx, req.ID)
	if err != nil {
		return articleparam.GetStorefrontArticleResponse{}, richerror.New(op).WithWarpError(err)
	}
	if !article.DeletedAt.IsZero() {
		return articleparam.GetStorefrontArticleResponse{}, richerror.New(op).WithKind(richerror.KindNotFound).WithMessageKey("global.not-found")
	}
	res, err := s.mediaSvc.GetURLByMediaID(ctx, mediaparam.GetUrlByIDRequest{ID: article.MediaID})
	if err != nil {
		return articleparam.GetStorefrontArticleResponse{}, richerror.New(op).WithWarpError(err)
	}
	a := articleparam.StorefrontArticleDetail{
		ID:          article.ID,
		Title:       article.Title,
		Slug:        article.Slug,
		Body:        article.Body,
		IsPublished: article.IsPublished,
		MediaUrl:    res.URL,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}
	return articleparam.GetStorefrontArticleResponse{
		Article: a,
	}, nil
}

func (s StorefrontService) GetAllStorefrontArticles(ctx context.Context, req articleparam.AllStorefrontArticleRequest) (articleparam.AllStorefrontArticleResponse, error) {
	const op = "ArticleService.GetAllStorefrontArticles"
	resp, err := s.repo.All(ctx, req.Filter, req.Pagination)
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
		return articleparam.AllStorefrontArticleResponse{}, err
	}

	for _, a := range resp.Articles {
		articlesInfo = append(articlesInfo, articleparam.StorefrontArticle{
			ID:          a.ID,
			Title:       a.Title,
			Slug:        a.Slug,
			IsPublished: a.IsPublished,
			MediaUrl:    urlMap.URLs[a.MediaID],
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
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
