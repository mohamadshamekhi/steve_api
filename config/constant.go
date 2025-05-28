package config

import "time"

const (
	defaultPrefix       = "ST_"
	defaultDelimiter    = "."
	defaultSeparator    = "__"
	defaultYamlFilePath = "config.yml"

	AccessTokenSubject                   = "ac"
	RefreshTokenSubject                  = "rt"
	BackOfficeAccessTokenExpireDuration  = time.Hour * 24 * 7
	BackOfficeRefreshTokenExpireDuration = time.Hour * 24 * 15
	AuthMiddlewareContextKey             = "claims"
	BcryptCost                           = 7
)
