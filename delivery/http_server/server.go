package httpserver

import (
	"fmt"
	"net/http"

	"github.com/abolfazlnorzad/steve/config"
	articlehandler "github.com/abolfazlnorzad/steve/delivery/http_server/article"
	attributehandler "github.com/abolfazlnorzad/steve/delivery/http_server/attribute"
	attributedefinitionhandler "github.com/abolfazlnorzad/steve/delivery/http_server/attributedefinition"
	backofficehandler "github.com/abolfazlnorzad/steve/delivery/http_server/backoffice"
	backofficeauthhandler "github.com/abolfazlnorzad/steve/delivery/http_server/backoffice/auth"
	brandhandler "github.com/abolfazlnorzad/steve/delivery/http_server/brand"
	carthandler "github.com/abolfazlnorzad/steve/delivery/http_server/cart"
	categoryhandler "github.com/abolfazlnorzad/steve/delivery/http_server/category"
	commenthandler "github.com/abolfazlnorzad/steve/delivery/http_server/comment"
	discounthandler "github.com/abolfazlnorzad/steve/delivery/http_server/discount"
	mediahandler "github.com/abolfazlnorzad/steve/delivery/http_server/media"
	middleware2 "github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	orderhandler "github.com/abolfazlnorzad/steve/delivery/http_server/order"
	producthandler "github.com/abolfazlnorzad/steve/delivery/http_server/product"
	repairhandler "github.com/abolfazlnorzad/steve/delivery/http_server/repair"
	sellhandler "github.com/abolfazlnorzad/steve/delivery/http_server/sell"
	shippinghandler "github.com/abolfazlnorzad/steve/delivery/http_server/shipping"
	userhandler "github.com/abolfazlnorzad/steve/delivery/http_server/user"
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config                config.Config
	Router                *echo.Echo
	backofficeAuthHandler backofficeauthhandler.Handler
	brandHandler          brandhandler.Handler
	mediaHandler          mediahandler.Handler
	categoryHandler       categoryhandler.Handler
	articleHandler        articlehandler.Handler
	attrHandler           attributehandler.Handler
	attrDefHandler        attributedefinitionhandler.Handler
	productHandler        producthandler.Handler
	backofficeHandler     backofficehandler.Handler
	userHandler           userhandler.Handler
	commentHandler        commenthandler.Handler
	discountHandler       discounthandler.Handler
	cartHandler           carthandler.Handler
	orderHandler          orderhandler.Handler
	shippingHandler       shippinghandler.Handler
	repairHandler         repairhandler.Handler
	sellHandler           sellhandler.Handler
}

func New(cfg config.Config, backofficeAuthHandler backofficeauthhandler.Handler, brandHandler brandhandler.Handler, mediaHandler mediahandler.Handler, categoryHandler categoryhandler.Handler,
	articleHandler articlehandler.Handler, attrHandler attributehandler.Handler, attrDefHandler attributedefinitionhandler.Handler,
	productHandler producthandler.Handler, backofficeHandler backofficehandler.Handler, userHandler userhandler.Handler,
	commentHandler commenthandler.Handler, discountHandler discounthandler.Handler, cartHandler carthandler.Handler,
	orderHandler orderhandler.Handler, shippingHandler shippinghandler.Handler, repairHandler repairhandler.Handler, sellHandler sellhandler.Handler,

) Server {
	return Server{
		config:                cfg,
		Router:                echo.New(),
		backofficeAuthHandler: backofficeAuthHandler,
		brandHandler:          brandHandler,
		mediaHandler:          mediaHandler,
		categoryHandler:       categoryHandler,
		articleHandler:        articleHandler,
		attrHandler:           attrHandler,
		attrDefHandler:        attrDefHandler,
		productHandler:        productHandler,
		backofficeHandler:     backofficeHandler,
		userHandler:           userHandler,
		commentHandler:        commentHandler,
		discountHandler:       discountHandler,
		cartHandler:           cartHandler,
		orderHandler:          orderHandler,
		shippingHandler:       shippingHandler,
		repairHandler:         repairHandler,
		sellHandler:           sellHandler,
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
			errMsg := ""
			if v.Error != nil {
				errMsg = v.Error.Error()
			}
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
				"error", errMsg,
			)

			return nil
		},
	}))

	s.Router.Use(middleware.Recover())

	//set routes
	s.backofficeAuthHandler.SetRoutes(s.Router)

	s.mediaHandler.SetRoutes(s.Router)
	s.brandHandler.SetRoutes(s.Router)
	s.categoryHandler.SetRoutes(s.Router)
	s.articleHandler.SetRoutes(s.Router)
	s.attrHandler.SetRoutes(s.Router)
	s.attrDefHandler.SetRoutes(s.Router)
	s.productHandler.SetRoutes(s.Router)
	s.backofficeHandler.SetRoutes(s.Router)
	s.userHandler.SetRoutes(s.Router)
	s.commentHandler.SetRoutes(s.Router)
	s.discountHandler.SetRoutes(s.Router)
	s.cartHandler.SetRoutes(s.Router)
	s.orderHandler.SetRoutes(s.Router)
	s.shippingHandler.SetRoutes(s.Router)
	s.repairHandler.SetRoutes(s.Router)
	s.sellHandler.SetRoutes(s.Router)
	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	fmt.Printf("start echo server on %s\n", address)
	if err := s.Router.Start(address); err != nil {
		fmt.Println("router start error", err)
	}
}
