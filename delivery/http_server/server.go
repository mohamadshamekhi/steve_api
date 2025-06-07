package httpserver

import (
	"fmt"
	"github.com/abolfazlnorzad/steve/config"
	backofficeauthhandler "github.com/abolfazlnorzad/steve/delivery/http_server/backoffice/auth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
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

func (s Server) Serve() {
	// todo -please fix cors policy configuration.
	s.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderAuthorization, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowOrigin, echo.HeaderAccessControlAllowMethods},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))
	s.Router.Use(middleware.RequestID())

	s.Router.Use(middleware.Logger())

	s.Router.Use(middleware.Recover())

	//set routes
	s.backofficeAuthHandler.SetRoutes(s.Router)

	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	fmt.Printf("start echo server on %s\n", address)
	if err := s.Router.Start(address); err != nil {
		fmt.Println("router start error", err)
	}
}
