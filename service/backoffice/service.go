package backofficeservice

import authservice "github.com/abolfazlnorzad/steve/service/auth"

type Config struct {
	Prefix string `koanf:"prefix"`
}
type Service struct {
	config  Config
	authSvc authservice.Service
}

func New(cfg Config, authSvc authservice.Service) Service {
	return Service{
		config:  cfg,
		authSvc: authSvc,
	}
}
