package http_server

import (
	"github.com/abolfazlnorzad/steve/config"
	backofficeauthhandler "github.com/abolfazlnorzad/steve/delivery/http_server/backoffice/auth"
	"github.com/labstack/echo/v4"
)

type Server struct {
	config                config.Config
	Router                *echo.Echo
	backofficeAuthHandler backofficeauthhandler.Handler
}

func New(cfg config.Config, backofficeAuthHandler backofficeauthhandler.Handler) Server {
	return Server{
		config:                cfg,
		Router:                echo.New(),
		backofficeAuthHandler: backofficeAuthHandler,
	}
}
