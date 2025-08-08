package mediaservice

import (
	"context"

	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s Service) DeleteMedia(ctx context.Context, req mediaparam.DeleteMediaRequest) (mediaparam.DeleteMediaResponse, error) {
	const op = "mediaservice.DeleteMedia"

	if err := s.provider.Delete(ctx, &req.Media); err != nil {

		logger.L().Error("warning: failed to delete physical file from storage: %v\n", err)
	}

	if err := s.repository.DeleteMedia(ctx, req.Media.ID); err != nil {
		return mediaparam.DeleteMediaResponse{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return mediaparam.DeleteMediaResponse{}, nil
}
