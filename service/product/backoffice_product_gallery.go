package productservice

import (
	"context"
	"fmt"
	"github.com/abolfazlnorzad/steve/entity"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	productparam "github.com/abolfazlnorzad/steve/param/product"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s BackofficeService) AddGallery(ctx context.Context, req productparam.AddProductGalleryRequest) (productparam.AddProductGalleryResponse, error) {
	const op = "productservice.AddGallery"

	g := entity.ProductGallery{
		MediaID:   req.MediaID,
		URL:       req.MediaURL,
		IsPrimary: req.IsPrimary,
	}
	if err := s.repo.AddProductGallery(ctx, g); err != nil {
		return productparam.AddProductGalleryResponse{}, richerror.New(op).WithWarpError(err)
	}

	return productparam.AddProductGalleryResponse{}, nil
}

func (s BackofficeService) DeleteGallery(ctx context.Context, req productparam.DeleteProductGalleryRequest) (productparam.DeleteProductGalleryResponse, error) {
	const op = "productservice.DeleteGallery"

	r, e := s.mediaSvc.FindMediaByID(ctx, mediaparam.GetMediaByIDRequest{ID: req.MediaID})
	if e != nil {
		return productparam.DeleteProductGalleryResponse{}, richerror.New(op).WithWarpError(e)
	}

	go func() {
		if _, err := s.mediaSvc.DeleteMedia(ctx, mediaparam.DeleteMediaRequest{Media: r.Media}); err != nil {
			logger.L().Error(fmt.Sprintf("can't delete product gallery media - media id :%s - err : %v"))
		}
	}()

	if err := s.repo.DeleteProductGallery(ctx, req.ProductID, req.MediaID); err != nil {
		return productparam.DeleteProductGalleryResponse{}, richerror.New(op).WithWarpError(err)
	}
	return productparam.DeleteProductGalleryResponse{}, nil

}
