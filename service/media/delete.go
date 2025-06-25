package mediaservice

import (
	"context"
	"errors"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s *Service) DeleteMedia(ctx context.Context, request mediaparam.DeleteMediaRequest) (mediaparam.DeleteMediaResponse, error) {
	const op = "mediaservice.DeleteMedia"

	media, err := s.repository.GetMediaByID(ctx, request.ID)
	if err != nil {
		var rErr *richerror.RichError
		if errors.As(err, &rErr) && rErr.Kind() == richerror.KindNotFound {
			return mediaparam.DeleteMediaResponse{}, err
		}
		return mediaparam.DeleteMediaResponse{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	if err := s.provider.Delete(ctx, &media); err != nil {

		logger.L().Error("warning: failed to delete physical file from storage: %v\n", err)
	}

	if err := s.repository.DeleteMedia(ctx, media.ID); err != nil {
		return mediaparam.DeleteMediaResponse{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}

	return mediaparam.DeleteMediaResponse{}, nil
}
