package authservice

import "time"

type Config struct {
	SignKey              string
	AccessSubject        string
	AccessExpirationTime time.Duration
}

type Service struct {
	cfg Config
}

func New(cfg Config) Service {
	return Service{cfg: cfg}
}
