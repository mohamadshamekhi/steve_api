package mediaparam

import "github.com/abolfazlnorzad/steve/entity"

type GetUrlRequest struct {
	Media entity.Media `json:"media"`
}

type GetUrlResponse struct {
	URL string `json:"url"`
}
