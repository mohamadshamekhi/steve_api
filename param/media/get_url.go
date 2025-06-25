package mediaparam

import "github.com/abolfazlnorzad/steve/entity"

type GetUrlRequest struct {
	Media entity.Media
}

type GetUrlResponse struct {
	URL string `json:"url"`
}
