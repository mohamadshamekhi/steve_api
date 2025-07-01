package mediaservice

import (
	"context"
	"errors"
	"fmt"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/abolfazlnorzad/steve/pkg/richerror"
	"github.com/redis/go-redis/v9"
)

func (s Service) GetURL(ctx context.Context, req mediaparam.GetUrlRequest) (mediaparam.GetUrlResponse, error) {
	const op = "media.service.GetURL"

	if req.Media.IsPublic {
		return mediaparam.GetUrlResponse{URL: fmt.Sprintf("%s/%s", s.config.PublicBaseURL, req.Media.Path)}, nil
	}

	cacheKey := fmt.Sprintf("media:url:%s", req.Media.ID)

	cachedURL, err := s.cache.Get(ctx, cacheKey)
	if err == nil {
		return mediaparam.GetUrlResponse{URL: cachedURL}, nil
	}

	if !errors.Is(err, redis.Nil) {
		logger.L().Error("cache provider error on GetURL, generating new url anyway",
			"error", err,
			"key", cacheKey)
	}

	newURL, err := s.provider.GetURL(ctx, &req.Media)
	if err != nil {
		return mediaparam.GetUrlResponse{}, richerror.New(op).WithWarpError(err).WithKind(richerror.KindUnexpected)
	}
	if cacheErr := s.cache.Set(ctx, cacheKey, newURL, s.config.PresignedPrivateURLExpiration); cacheErr != nil {
		logger.L().Error("failed to cache media url", "error", cacheErr, "key", cacheKey)
	}
	return mediaparam.GetUrlResponse{URL: newURL}, nil
}
