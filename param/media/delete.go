package mediaparam

import "github.com/abolfazlnorzad/steve/entity"

type DeleteMediaRequest struct {
	Media entity.Media `json:"media"`
}

type DeleteMediaResponse struct {
}
