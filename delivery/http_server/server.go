package httpserver

import (
	"fmt"
	"github.com/abolfazlnorzad/steve/config"
	backofficeauthhandler "github.com/abolfazlnorzad/steve/delivery/http_server/backoffice/auth"
	brandhandler "github.com/abolfazlnorzad/steve/delivery/http_server/brand"
	mediahandler "github.com/abolfazlnorzad/steve/delivery/http_server/media"
	middleware2 "github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type Server struct {
	config                config.Config
	Router                *echo.Echo
	backofficeAuthHandler backofficeauthhandler.Handler
	brandHandler          brandhandler.Handler
	mediaHandler          mediahandler.Handler
}

func New(cfg config.Config, backofficeAuthHandler backofficeauthhandler.Handler, brandHandler brandhandler.Handler, mediaHandler mediahandler.Handler) Server {
	return Server{
		config:                cfg,
		Router:                echo.New(),
		backofficeAuthHandler: backofficeAuthHandler,
		brandHandler:          brandHandler,
		mediaHandler:          mediaHandler,
	}
}

func (s Server) Serve() {
	s.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderAuthorization, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowOrigin, echo.HeaderAccessControlAllowMethods},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	translator, _ := i18n.New("config/locales", "fa")

	s.Router.HTTPErrorHandler = middleware2.NewCentralErrorHandler(translator)

	s.Router.Use(middleware.RequestID())

	s.Router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:           true,
		LogStatus:        true,
		LogHost:          true,
		LogRemoteIP:      true,
		LogRequestID:     true,
		LogMethod:        true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogLatency:       true,
		LogError:         true,
		LogProtocol:      true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.L().Info("request",
					"request_id", v.RequestID,
					"host", v.Host,
					"content_length", v.ContentLength,
					"protocol", v.Protocol,
					"method", v.Method,
					"latency", v.Latency,
					"error", v.Error,
					"remote_ip", v.RemoteIP,
					"response_size", v.ResponseSize,
					"uri", v.URI,
					"status", v.Status,
				)
			}

			return nil
		},
	}))

	s.Router.Use(middleware.Recover())

	//set routes
	s.backofficeAuthHandler.SetRoutes(s.Router)

	s.mediaHandler.SetRoutes(s.Router)
	s.brandHandler.SetRoutes(s.Router)

	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	fmt.Printf("start echo server on %s\n", address)
	if err := s.Router.Start(address); err != nil {
		fmt.Println("router start error", err)
	}
}
