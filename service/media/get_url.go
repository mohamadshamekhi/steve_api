package mediaservice

import (
	"context"
	"fmt"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s *Service) GetURL(ctx context.Context, req mediaparam.GetUrlRequest) (mediaparam.GetUrlResponse, error) {
	const op = "media.service.GetURL"

	if req.Media.IsPublic {
		return mediaparam.GetUrlResponse{URL: fmt.Sprintf("%s/%s", s.config.PublicBaseURL, req.Media.Path)}, nil
	}

	url, err := s.provider.GetURL(ctx, req.Media.Path, s.config.PresignedPrivateURLExpiration)
	if err != nil {
		return mediaparam.GetUrlResponse{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	return mediaparam.GetUrlResponse{URL: url}, nil
}
