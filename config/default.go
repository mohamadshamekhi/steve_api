package config

import authservice "github.com/abolfazlnorzad/steve/service/auth"

func Default() Config {
	return Config{
		HTTPServer: HTTPServer{
			Port: 7777,
		},
		Auth: authservice.Config{
			AccessExpirationTime: BackOfficeAccessTokenExpireDuration,
			AccessSubject:        AccessTokenSubject,
		},
	}
}
