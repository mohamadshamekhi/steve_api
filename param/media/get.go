package mediaparam

import "github.com/abolfazlnorzad/steve/entity"

type GetMediaByIDRequest struct {
	ID string `json:"id"`
}

type GetMediaByIDResponse struct {
	Media entity.Media `json:"media"`
}

type GetUrlsByIDsRequest struct {
	IDs []string `json:"ids"`
}
type GetUrlsByIDsResponse struct {
	URLs map[string]string `json:"urls"`
}
