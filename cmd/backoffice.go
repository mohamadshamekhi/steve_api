package main

import (
	natsadapter "github.com/abolfazlnorzad/steve/adapter/nats"
	redisadapter "github.com/abolfazlnorzad/steve/adapter/redis"
	"github.com/abolfazlnorzad/steve/config"
	httpserver "github.com/abolfazlnorzad/steve/delivery/http_server"
	backofficeauthhandler "github.com/abolfazlnorzad/steve/delivery/http_server/backoffice/auth"
	brandhandler "github.com/abolfazlnorzad/steve/delivery/http_server/brand"
	"github.com/abolfazlnorzad/steve/internal/seeder"
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	mongodbpermission "github.com/abolfazlnorzad/steve/repository/mongo/acl/permission"
	mongodbrole "github.com/abolfazlnorzad/steve/repository/mongo/acl/role"
	mongodbbackoffice "github.com/abolfazlnorzad/steve/repository/mongo/backoffice"
	mongodbbrand "github.com/abolfazlnorzad/steve/repository/mongo/brand"
	redisotp "github.com/abolfazlnorzad/steve/repository/redis/otp"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	backofficeservice "github.com/abolfazlnorzad/steve/service/backoffice"
	brandservice "github.com/abolfazlnorzad/steve/service/brand"
	backofficevalidatorauth "github.com/abolfazlnorzad/steve/validator/backoffice/auth"
	brandvalidator "github.com/abolfazlnorzad/steve/validator/brand"
)

func main() {
	cfg := config.C()
	redisRepo := redisadapter.New(cfg.Redis)
	nats := natsadapter.New(natsadapter.Config{
		Host: "localhost",
		Port: 4222,
	})
	authSvc := authservice.New(cfg.Auth, redisotp.New(redisRepo), nats)
	mongoRepo := mongodb.NewMongoDB(cfg.Mongo)
	backOfficeMongoRepo := mongodbbackoffice.New(mongoRepo)
	backOfficeSvc := backofficeservice.New(cfg.BackOffice, authSvc, backOfficeMongoRepo)
	backOfficeVld := backofficevalidatorauth.New(backOfficeMongoRepo)

	translator, _ := i18n.New("config/locales", "fa")

	permissionRepo := mongodbpermission.New(mongoRepo)
	roleRepo := mongodbrole.New(mongoRepo)
	seeder.Run(permissionRepo, roleRepo, backOfficeMongoRepo)

	brandRepo := mongodbbrand.New(mongoRepo)
	backofficebrandsvc := brandservice.NewBackofficeService(brandRepo)

	brandVld := brandvalidator.New(brandRepo)
	backOfficeHandler := backofficeauthhandler.New(backOfficeSvc, backOfficeVld, translator)
	aclSvc := aclservice.New(roleRepo)
	brandHandler := brandhandler.New(backofficebrandsvc, translator, brandVld, authSvc, cfg.Auth, aclSvc)
	server := httpserver.New(cfg, backOfficeHandler, brandHandler)
	server.Serve()
}
