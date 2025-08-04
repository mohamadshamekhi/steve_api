package backofficeparam

import "time"

type LogoutRequest struct {
	JWTID    string    `json:"jwt_id"`
	ExpireAt time.Time `json:"expire_at"`
	Kind     string    `json:"kind"`
}

type LogoutResponse struct {
}
