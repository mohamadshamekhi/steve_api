package mediaservice

import (
	"context"
	mediaparam "github.com/abolfazlnorzad/steve/param/media"
)

func (s Service) FindMediaByID(ctx context.Context, req mediaparam.GetMediaByIDRequest) (mediaparam.GetMediaByIDResponse, error) {
	const op = "mediaservice.FindMediaByID"
	media, err := s.repository.GetMediaByID(ctx, req.ID)
	if err != nil {
		return mediaparam.GetMediaByIDResponse{}, err
	}
	return mediaparam.GetMediaByIDResponse{Media: media}, nil
}
