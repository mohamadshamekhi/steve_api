package main

import (
	natsadapter "github.com/abolfazlnorzad/steve/adapter/nats"
	redisadapter "github.com/abolfazlnorzad/steve/adapter/redis"
	liarastorage "github.com/abolfazlnorzad/steve/adapter/storage/liara"
	"github.com/abolfazlnorzad/steve/config"
	httpserver "github.com/abolfazlnorzad/steve/delivery/http_server"
	articlehandler "github.com/abolfazlnorzad/steve/delivery/http_server/article"
	attributehandler "github.com/abolfazlnorzad/steve/delivery/http_server/attribute"
	attributedefinitionhandler "github.com/abolfazlnorzad/steve/delivery/http_server/attributedefinition"
	backofficehandler "github.com/abolfazlnorzad/steve/delivery/http_server/backoffice"
	backofficeauthhandler "github.com/abolfazlnorzad/steve/delivery/http_server/backoffice/auth"
	brandhandler "github.com/abolfazlnorzad/steve/delivery/http_server/brand"
	categoryhandler "github.com/abolfazlnorzad/steve/delivery/http_server/category"
	mediahandler "github.com/abolfazlnorzad/steve/delivery/http_server/media"
	producthandler "github.com/abolfazlnorzad/steve/delivery/http_server/product"
	userhandler "github.com/abolfazlnorzad/steve/delivery/http_server/user"
	"github.com/abolfazlnorzad/steve/internal/seeder"
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	mongodbpermission "github.com/abolfazlnorzad/steve/repository/mongo/acl/permission"
	mongodbrole "github.com/abolfazlnorzad/steve/repository/mongo/acl/role"
	mongodbarticle "github.com/abolfazlnorzad/steve/repository/mongo/article"
	mongodbattribute "github.com/abolfazlnorzad/steve/repository/mongo/attribute"
	mongodbattributedefinition "github.com/abolfazlnorzad/steve/repository/mongo/attribute_definition"
	mongodbbackoffice "github.com/abolfazlnorzad/steve/repository/mongo/backoffice"
	mongodbbrand "github.com/abolfazlnorzad/steve/repository/mongo/brand"
	mongodbcategory "github.com/abolfazlnorzad/steve/repository/mongo/category"
	mongodbmedia "github.com/abolfazlnorzad/steve/repository/mongo/media"
	mongodbproduct "github.com/abolfazlnorzad/steve/repository/mongo/product"
	mongodbuser "github.com/abolfazlnorzad/steve/repository/mongo/user"
	redisotp "github.com/abolfazlnorzad/steve/repository/redis/otp"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	articleservice "github.com/abolfazlnorzad/steve/service/article"
	attributeservice "github.com/abolfazlnorzad/steve/service/attribute"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	backofficeservice "github.com/abolfazlnorzad/steve/service/backoffice"
	brandservice "github.com/abolfazlnorzad/steve/service/brand"
	categoryservice "github.com/abolfazlnorzad/steve/service/category"
	mediaservice "github.com/abolfazlnorzad/steve/service/media"
	productservice "github.com/abolfazlnorzad/steve/service/product"
	userservice "github.com/abolfazlnorzad/steve/service/user"
	articlevalidator "github.com/abolfazlnorzad/steve/validator/article"
	attributevalidator "github.com/abolfazlnorzad/steve/validator/attribute"
	attributedefinitionvalidator "github.com/abolfazlnorzad/steve/validator/attributedefinition"
	backofficevalidator "github.com/abolfazlnorzad/steve/validator/backoffice"
	backofficevalidatorauth "github.com/abolfazlnorzad/steve/validator/backoffice/auth"
	brandvalidator "github.com/abolfazlnorzad/steve/validator/brand"
	categoryvalidator "github.com/abolfazlnorzad/steve/validator/category"
	mediavalidator "github.com/abolfazlnorzad/steve/validator/media"
	productvalidator "github.com/abolfazlnorzad/steve/validator/product"
	uservalidator "github.com/abolfazlnorzad/steve/validator/user"
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
	attrRepo := mongodbattribute.New(mongoRepo)
	attrDefRepo := mongodbattributedefinition.New(mongoRepo)

	attrVld := attributevalidator.New(attrRepo)
	attrSvc := attributeservice.NewAttributeService(attrRepo, attrDefRepo)
	categoryVld := categoryvalidator.New(categoryRepo, mediaSvc, attrSvc)
	categorySvc := categoryservice.NewBackofficeService(categoryRepo, mediaSvc, attrSvc)
	storefrontCategorySvc := categoryservice.NewStorefrontService(categoryRepo, mediaSvc)
	categoryHandler := categoryhandler.New(categorySvc, categoryVld, translator, authSvc, cfg.Auth, aclSvc, storefrontCategorySvc)
	articleRepo := mongodbarticle.New(mongoRepo)
	articleVld := articlevalidator.New(articleRepo, mediaSvc)
	articleSvc := articleservice.NewBackofficeService(articleRepo, mediaSvc)
	storefrontSvc := articleservice.NewStorefrontService(articleRepo, mediaSvc)
	storefrontVld := articlevalidator.NewStorefrontValidator(articleRepo)
	articleHandler := articlehandler.New(translator, authSvc, cfg.Auth, aclSvc, articleSvc, articleVld, storefrontSvc, storefrontVld)

	attrHandler := attributehandler.New(authSvc, cfg.Auth, aclSvc, translator, attrVld, attrSvc)

	attrDefVld := attributedefinitionvalidator.New(attrDefRepo)
	attrDefSvc := attributeservice.NewAttributeDefinition(attrDefRepo)
	attrDefHandler := attributedefinitionhandler.New(authSvc, cfg.Auth, aclSvc, translator, attrDefVld, attrDefSvc)

	productRepo := mongodbproduct.New(mongoRepo)
	productVld := productvalidator.New(categorySvc, backofficebrandsvc, productRepo, mediaSvc)
	productSvc := productservice.NewBackofficeService(productRepo, mediaSvc, categorySvc, backofficebrandsvc)

	productHandler := producthandler.NewHandler(translator, authSvc, cfg.Auth, aclSvc, productSvc, productVld)
	backofficeValidation := backofficevalidator.New(backOfficeMongoRepo)
	backofficeHandler := backofficehandler.New(backOfficeSvc, backofficeValidation, translator, authSvc, cfg.Auth, aclSvc)
	userRepo := mongodbuser.New(mongoRepo)
	userVld := uservalidator.New(userRepo)
	userSvc := userservice.NewStoreFrontService(cfg.User, userRepo, authSvc)
	userHandler := userhandler.New(userSvc, userVld, translator, authSvc, cfg.Auth)

	server := httpserver.New(cfg, backOfficeHandler, brandHandler, mediaHandler, categoryHandler, articleHandler, attrHandler, attrDefHandler, productHandler, backofficeHandler, userHandler)
	server.Serve()
}
