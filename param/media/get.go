package mediaparam

import "github.com/abolfazlnorzad/steve/entity"

type GetMediaByIDRequest struct {
	ID string `json:"id"`
}

type GetMediaByIDResponse struct {
	Media entity.Media `json:"media"`
}
