package main

import (
	natsadapter "github.com/abolfazlnorzad/steve/adapter/nats"
	redisadapter "github.com/abolfazlnorzad/steve/adapter/redis"
	liarastorage "github.com/abolfazlnorzad/steve/adapter/storage/liara"
	"github.com/abolfazlnorzad/steve/config"
	httpserver "github.com/abolfazlnorzad/steve/delivery/http_server"
	articlehandler "github.com/abolfazlnorzad/steve/delivery/http_server/article"
	backofficeauthhandler "github.com/abolfazlnorzad/steve/delivery/http_server/backoffice/auth"
	brandhandler "github.com/abolfazlnorzad/steve/delivery/http_server/brand"
	categoryhandler "github.com/abolfazlnorzad/steve/delivery/http_server/category"
	mediahandler "github.com/abolfazlnorzad/steve/delivery/http_server/media"
	"github.com/abolfazlnorzad/steve/internal/seeder"
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	mongodbpermission "github.com/abolfazlnorzad/steve/repository/mongo/acl/permission"
	mongodbrole "github.com/abolfazlnorzad/steve/repository/mongo/acl/role"
	mongodbarticle "github.com/abolfazlnorzad/steve/repository/mongo/article"
	mongodbbackoffice "github.com/abolfazlnorzad/steve/repository/mongo/backoffice"
	mongodbbrand "github.com/abolfazlnorzad/steve/repository/mongo/brand"
	mongodbcategory "github.com/abolfazlnorzad/steve/repository/mongo/category"
	mongodbmedia "github.com/abolfazlnorzad/steve/repository/mongo/media"
	redisotp "github.com/abolfazlnorzad/steve/repository/redis/otp"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	articleservice "github.com/abolfazlnorzad/steve/service/article"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	backofficeservice "github.com/abolfazlnorzad/steve/service/backoffice"
	brandservice "github.com/abolfazlnorzad/steve/service/brand"
	categoryservice "github.com/abolfazlnorzad/steve/service/category"
	mediaservice "github.com/abolfazlnorzad/steve/service/media"
	articlevalidator "github.com/abolfazlnorzad/steve/validator/article"
	backofficevalidatorauth "github.com/abolfazlnorzad/steve/validator/backoffice/auth"
	brandvalidator "github.com/abolfazlnorzad/steve/validator/brand"
	categoryvalidator "github.com/abolfazlnorzad/steve/validator/category"
	mediavalidator "github.com/abolfazlnorzad/steve/validator/media"
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
	mediaRepo := mongodbmedia.New(mongoRepo)

	l := liarastorage.New(cfg.ObjectStorage, cfg.Media)
	mediaSvc := mediaservice.New(l, mediaRepo, cfg.Media, redisRepo)
	backofficebrandsvc := brandservice.NewBackofficeService(brandRepo, mediaSvc)

	brandVld := brandvalidator.New(brandRepo, mediaSvc)
	backOfficeHandler := backofficeauthhandler.New(backOfficeSvc, backOfficeVld, translator)
	aclSvc := aclservice.New(roleRepo)
	brandStorefrontSvc := brandservice.NewStorefrontService(brandRepo, mediaSvc)
	brandHandler := brandhandler.New(backofficebrandsvc, translator, brandVld, authSvc, cfg.Auth, aclSvc, brandStorefrontSvc)
	mediaVld := mediavalidator.New(mediaRepo)

	mediaHandler := mediahandler.New(mediaVld, mediaSvc, translator, authSvc, cfg.Auth)
	categoryRepo := mongodbcategory.New(mongoRepo)
	categoryVld := categoryvalidator.New(categoryRepo, mediaSvc)
	categorySvc := categoryservice.NewBackofficeService(categoryRepo, mediaSvc)
	categoryHandler := categoryhandler.New(categorySvc, categoryVld, translator, authSvc, cfg.Auth, aclSvc)
	articleRepo := mongodbarticle.New(mongoRepo)
	articleVld := articlevalidator.New(articleRepo, mediaSvc)
	articleSvc := articleservice.NewBackofficeService(articleRepo, mediaSvc)
	articleHandler := articlehandler.New(translator, authSvc, cfg.Auth, aclSvc, articleSvc, articleVld)
	server := httpserver.New(cfg, backOfficeHandler, brandHandler, mediaHandler, categoryHandler, articleHandler)
	server.Serve()
}
