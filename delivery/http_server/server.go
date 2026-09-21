package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abolfazlnorzad/steve/config"
	accountinghandler "github.com/abolfazlnorzad/steve/delivery/http_server/accounting"
	aclhandler "github.com/abolfazlnorzad/steve/delivery/http_server/acl/handler"
	articlehandler "github.com/abolfazlnorzad/steve/delivery/http_server/article"
	attributehandler "github.com/abolfazlnorzad/steve/delivery/http_server/attribute"
	attributedefinitionhandler "github.com/abolfazlnorzad/steve/delivery/http_server/attributedefinition"
	backofficehandler "github.com/abolfazlnorzad/steve/delivery/http_server/backoffice"
	backofficeauthhandler "github.com/abolfazlnorzad/steve/delivery/http_server/backoffice/auth"
	bannerhandler "github.com/abolfazlnorzad/steve/delivery/http_server/banner"
	brandhandler "github.com/abolfazlnorzad/steve/delivery/http_server/brand"
	carthandler "github.com/abolfazlnorzad/steve/delivery/http_server/cart"
	categoryhandler "github.com/abolfazlnorzad/steve/delivery/http_server/category"
	commenthandler "github.com/abolfazlnorzad/steve/delivery/http_server/comment"
	contactushandler "github.com/abolfazlnorzad/steve/delivery/http_server/contactus"
	discounthandler "github.com/abolfazlnorzad/steve/delivery/http_server/discount"
	inventoryhandler "github.com/abolfazlnorzad/steve/delivery/http_server/inventory"
	mediahandler "github.com/abolfazlnorzad/steve/delivery/http_server/media"
	middleware2 "github.com/abolfazlnorzad/steve/delivery/http_server/middleware"
	orderhandler "github.com/abolfazlnorzad/steve/delivery/http_server/order"
	orderdeliveryhandler "github.com/abolfazlnorzad/steve/delivery/http_server/orderdelivery"
	partnerhandler "github.com/abolfazlnorzad/steve/delivery/http_server/partner"
	paymenthandler "github.com/abolfazlnorzad/steve/delivery/http_server/payment"
	producthandler "github.com/abolfazlnorzad/steve/delivery/http_server/product"
	purchasehandler "github.com/abolfazlnorzad/steve/delivery/http_server/purchase"
	ratinghandler "github.com/abolfazlnorzad/steve/delivery/http_server/rating"
	referralhandler "github.com/abolfazlnorzad/steve/delivery/http_server/referral"
	repairhandler "github.com/abolfazlnorzad/steve/delivery/http_server/repair"
	salesreturnhandler "github.com/abolfazlnorzad/steve/delivery/http_server/salesreturn"
	searchhandler "github.com/abolfazlnorzad/steve/delivery/http_server/search"
	sellhandler "github.com/abolfazlnorzad/steve/delivery/http_server/sell"
	shippinghandler "github.com/abolfazlnorzad/steve/delivery/http_server/shipping"
	supplierhandler "github.com/abolfazlnorzad/steve/delivery/http_server/supplier"
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
	aclHandler            aclhandler.Handler
	contactUsHandler      contactushandler.Handler
	ratingHandler         ratinghandler.Handler
	bannerHandler         bannerhandler.Handler
	paymentHandler        paymenthandler.Handler
	searchHandler         searchhandler.Handler
	partnerHandler        partnerhandler.Handler
	referralHandler       referralhandler.Handler
	supplierHandler       supplierhandler.Handler
	purchaseHandler       purchasehandler.Handler
	inventoryHandler      inventoryhandler.Handler
	accountingHandler     accountinghandler.Handler
	salesReturnHandler    salesreturnhandler.Handler
	orderDeliveryHandler  orderdeliveryhandler.Handler
}

func New(cfg config.Config, backofficeAuthHandler backofficeauthhandler.Handler, brandHandler brandhandler.Handler, mediaHandler mediahandler.Handler, categoryHandler categoryhandler.Handler,
	articleHandler articlehandler.Handler, attrHandler attributehandler.Handler, attrDefHandler attributedefinitionhandler.Handler,
	productHandler producthandler.Handler, backofficeHandler backofficehandler.Handler, userHandler userhandler.Handler,
	commentHandler commenthandler.Handler, discountHandler discounthandler.Handler, cartHandler carthandler.Handler,
	orderHandler orderhandler.Handler, shippingHandler shippinghandler.Handler, repairHandler repairhandler.Handler, sellHandler sellhandler.Handler,
	aclHandler aclhandler.Handler, contactUsHandler contactushandler.Handler, ratingHandler ratinghandler.Handler,
	bannerHandler bannerhandler.Handler, paymentHandler paymenthandler.Handler, searchHandler searchhandler.Handler,
	partnerHandler partnerhandler.Handler, referralHandler referralhandler.Handler,
	supplierHandler supplierhandler.Handler, purchaseHandler purchasehandler.Handler,
	inventoryHandler inventoryhandler.Handler,
	accountingHandler accountinghandler.Handler,
	salesReturnHandler salesreturnhandler.Handler,
	orderDeliveryHandler orderdeliveryhandler.Handler,
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
		aclHandler:            aclHandler,
		contactUsHandler:      contactUsHandler,
		ratingHandler:         ratingHandler,
		bannerHandler:         bannerHandler,
		paymentHandler:        paymentHandler,
		searchHandler:         searchHandler,
		partnerHandler:        partnerHandler,
		referralHandler:       referralHandler,
		supplierHandler:       supplierHandler,
		purchaseHandler:       purchaseHandler,
		inventoryHandler:      inventoryHandler,
		accountingHandler:     accountingHandler,
		salesReturnHandler:    salesReturnHandler,
		orderDeliveryHandler:  orderDeliveryHandler,
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
					"error_msg", errMsg,
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
	s.aclHandler.SetRoutes(s.Router)
	s.contactUsHandler.SetRoutes(s.Router)
	s.ratingHandler.SetRoutes(s.Router)
	s.bannerHandler.SetRoutes(s.Router)
	s.paymentHandler.SetRoutes(s.Router)
	s.searchHandler.SetRoutes(s.Router)
	s.partnerHandler.SetRoutes(s.Router)
	s.referralHandler.SetRoutes(s.Router)
	s.supplierHandler.SetRoutes(s.Router)
	s.purchaseHandler.SetRoutes(s.Router)
	s.inventoryHandler.SetRoutes(s.Router)
	s.accountingHandler.SetRoutes(s.Router)
	s.salesReturnHandler.SetRoutes(s.Router)
	s.orderDeliveryHandler.SetRoutes(s.Router)
	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	fmt.Printf("start echo server on %s\n", address)
	if err := s.Router.Start(address); err != nil {
		fmt.Println("router start error", err)
	}
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.Router.Shutdown(ctx)
}
