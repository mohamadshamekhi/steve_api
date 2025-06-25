package mediaservice

import (
	"context"
	"fmt"
	"github.com/abolfazlnorzad/steve/entity"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
)

func (s *Service) GetURL(ctx context.Context, media *entity.Media) (string, error) {
	const op = "media.service.GetURL"

	if media.IsPublic {
		return fmt.Sprintf("%s/%s", s.config.PublicBaseURL, media.Path), nil
	}

	url, err := s.provider.GetURL(ctx, media.Path, s.config.PresignedPrivateURLExpiration)
	if err != nil {
		return "", richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	return url, nil
}
