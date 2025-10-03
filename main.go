package main

import (
	natsadapter "github.com/abolfazlnorzad/steve/adapter/nats"
	redisadapter "github.com/abolfazlnorzad/steve/adapter/redis"
	liarastorage "github.com/abolfazlnorzad/steve/adapter/storage/liara"
	"github.com/abolfazlnorzad/steve/config"
	httpserver "github.com/abolfazlnorzad/steve/delivery/http_server"
	aclhandler "github.com/abolfazlnorzad/steve/delivery/http_server/acl/handler"
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
	orderhandler "github.com/abolfazlnorzad/steve/delivery/http_server/order"
	producthandler "github.com/abolfazlnorzad/steve/delivery/http_server/product"
	repairhandler "github.com/abolfazlnorzad/steve/delivery/http_server/repair"
	sellhandler "github.com/abolfazlnorzad/steve/delivery/http_server/sell"
	shippinghandler "github.com/abolfazlnorzad/steve/delivery/http_server/shipping"
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
	mongodbcart "github.com/abolfazlnorzad/steve/repository/mongo/cart"
	mongodbcategory "github.com/abolfazlnorzad/steve/repository/mongo/category"
	mongodbcomment "github.com/abolfazlnorzad/steve/repository/mongo/comment"
	mongodbdiscount "github.com/abolfazlnorzad/steve/repository/mongo/discount"
	mongodbmedia "github.com/abolfazlnorzad/steve/repository/mongo/media"
	mongodborder "github.com/abolfazlnorzad/steve/repository/mongo/order"
	mongodbproduct "github.com/abolfazlnorzad/steve/repository/mongo/product"
	mongodbrepair "github.com/abolfazlnorzad/steve/repository/mongo/repair"
	mongodbsell "github.com/abolfazlnorzad/steve/repository/mongo/sell"
	mongodbshipping "github.com/abolfazlnorzad/steve/repository/mongo/shipping"
	mongodbuser "github.com/abolfazlnorzad/steve/repository/mongo/user"
	redisotp "github.com/abolfazlnorzad/steve/repository/redis/otp"
	aclservice "github.com/abolfazlnorzad/steve/service/acl"
	permissionservice "github.com/abolfazlnorzad/steve/service/acl/permission"
	roleservice "github.com/abolfazlnorzad/steve/service/acl/role"
	articleservice "github.com/abolfazlnorzad/steve/service/article"
	attributeservice "github.com/abolfazlnorzad/steve/service/attribute"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	backofficeservice "github.com/abolfazlnorzad/steve/service/backoffice"
	brandservice "github.com/abolfazlnorzad/steve/service/brand"
	cartservice "github.com/abolfazlnorzad/steve/service/cart"
	categoryservice "github.com/abolfazlnorzad/steve/service/category"
	commentservice "github.com/abolfazlnorzad/steve/service/comment"
	discountservice "github.com/abolfazlnorzad/steve/service/discount"
	fetcherservice "github.com/abolfazlnorzad/steve/service/fetcher"
	mediaservice "github.com/abolfazlnorzad/steve/service/media"
	orderservice "github.com/abolfazlnorzad/steve/service/order"
	productservice "github.com/abolfazlnorzad/steve/service/product"
	repairservice "github.com/abolfazlnorzad/steve/service/repair"
	sellservice "github.com/abolfazlnorzad/steve/service/sell"
	settingservice "github.com/abolfazlnorzad/steve/service/setting"
	userservice "github.com/abolfazlnorzad/steve/service/user"
	articlevalidator "github.com/abolfazlnorzad/steve/validator/article"
	attributevalidator "github.com/abolfazlnorzad/steve/validator/attribute"
	attributedefinitionvalidator "github.com/abolfazlnorzad/steve/validator/attributedefinition"
	backofficevalidator "github.com/abolfazlnorzad/steve/validator/backoffice"
	backofficevalidatorauth "github.com/abolfazlnorzad/steve/validator/backoffice/auth"
	brandvalidator "github.com/abolfazlnorzad/steve/validator/brand"
	cartvalidator "github.com/abolfazlnorzad/steve/validator/cart"
	categoryvalidator "github.com/abolfazlnorzad/steve/validator/category"
	commentvalidator "github.com/abolfazlnorzad/steve/validator/comment"
	discountvalidator "github.com/abolfazlnorzad/steve/validator/discount"
	mediavalidator "github.com/abolfazlnorzad/steve/validator/media"
	ordervalidator "github.com/abolfazlnorzad/steve/validator/order"
	productvalidator "github.com/abolfazlnorzad/steve/validator/product"
	repairvalidator "github.com/abolfazlnorzad/steve/validator/repair"
	sellvalidator "github.com/abolfazlnorzad/steve/validator/sell"
	settingvalidator "github.com/abolfazlnorzad/steve/validator/setting"
	uservalidator "github.com/abolfazlnorzad/steve/validator/user"
)

func main() {
	cfg := config.C()
	redisRepo := redisadapter.New(cfg.Redis)
	nats := natsadapter.New(natsadapter.Config{
		Host: "localhost",
		Port: 4222,
	})
	redisotp := redisotp.New(redisRepo)
	authSvc := authservice.New(cfg.Auth, redisotp, nats)
	mongoRepo := mongodb.NewMongoDB(cfg.Mongo)
	backOfficeMongoRepo := mongodbbackoffice.New(mongoRepo)
	backOfficeSvc := backofficeservice.New(cfg.BackOffice, authSvc, backOfficeMongoRepo, redisotp)
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
	aclSvc := aclservice.New(roleRepo, backOfficeMongoRepo)
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
	productStorefrontSvc := productservice.NewStorefrontService(productRepo, mediaSvc, categorySvc, backofficebrandsvc)

	productHandler := producthandler.NewHandler(translator, authSvc, cfg.Auth, aclSvc, productSvc, productStorefrontSvc, productVld)
	backofficeValidation := backofficevalidator.New(backOfficeMongoRepo)
	backofficeHandler := backofficehandler.New(backOfficeSvc, backofficeValidation, translator, authSvc, cfg.Auth, aclSvc)
	userRepo := mongodbuser.New(mongoRepo)
	userVld := uservalidator.New(userRepo)
	userSvc := userservice.NewStoreFrontService(cfg.User, userRepo, authSvc, redisotp)
	backofficeUserSvc := userservice.NewBackofficeService(userRepo)
	userHandler := userhandler.New(backofficeUserSvc, userSvc, userVld, translator, authSvc, cfg.Auth, aclSvc)
	commentRepo := mongodbcomment.New(mongoRepo)
	commentVld := commentvalidator.New(commentRepo)

	commentableFetecher := fetcherservice.NewCommentableFetcher(*productRepo, *articleRepo)
	commenterFetecher := fetcherservice.NewCommenterFetcher(*userRepo, *backOfficeMongoRepo)
	commentSvc := commentservice.New(commentRepo, commenterFetecher, commentableFetecher)
	storefrontCommentSvc := commentservice.NewStorefrontService(commentRepo, commentableFetecher, commenterFetecher)
	commentHandler := commenthandler.New(commentSvc, translator, commentVld, authSvc, cfg.Auth, aclSvc, storefrontCommentSvc)
	discountRepo := mongodbdiscount.New(mongoRepo)
	discountVld := discountvalidator.New(discountRepo, categoryRepo, productRepo, userRepo)

	cartRepo := mongodbcart.New(mongoRepo)
	cartVld := cartvalidator.New(productRepo, cartRepo)
	cartSvc := cartservice.NewStorefrontService(cartRepo, productRepo, userRepo)
	discountbackofficeSvc := discountservice.NewBackofficeService(discountRepo, cartSvc, productSvc)
	discountHandler := discounthandler.New(discountbackofficeSvc, translator, discountVld, authSvc, cfg.Auth, aclSvc)
	cartBackofficeSvc := cartservice.NewBackofficeService(cartRepo, userRepo, productRepo)
	orderRepo := mongodborder.New(mongoRepo)
	orderVld := ordervalidator.New(discountbackofficeSvc)
	shippingRepo := mongodbshipping.New(mongoRepo)
	shippingVld := settingvalidator.NewShippingValidator(shippingRepo)
	shippingBakcofficeSvc := settingservice.NewBackofficeShippingService(shippingRepo)
	orderSvc := orderservice.New(cartSvc, discountbackofficeSvc, productSvc, orderRepo, userSvc, shippingBakcofficeSvc)
	cartHandler := carthandler.New(cartSvc, cartBackofficeSvc, orderSvc, cartVld, translator, authSvc, cfg.Auth, aclSvc)

	orderBackofficeSvc := orderservice.NewBackofficeService(orderRepo, backofficeUserSvc)
	orderHandler := orderhandler.New(orderSvc, orderBackofficeSvc, orderVld, translator, authSvc, cfg.Auth, aclSvc)

	shippingHandler := shippinghandler.New(shippingBakcofficeSvc, shippingVld, translator, authSvc, cfg.Auth, aclSvc)

	repairRepo := mongodbrepair.New(mongoRepo)
	repairVld := repairvalidator.New()
	repairStorfrontSvc := repairservice.NewStorefrontService(repairRepo, backofficeUserSvc, backofficebrandsvc, productSvc)
	repairBackofficeSvc := repairservice.NewBackofficeService(repairRepo, backofficeUserSvc, backofficebrandsvc, productSvc)
	reapirHandler := repairhandler.New(repairBackofficeSvc, repairStorfrontSvc, translator, repairVld, authSvc, cfg.Auth, aclSvc)

	sellRepo := mongodbsell.New(mongoRepo)
	sellVld := sellvalidator.New()
	sellStorefrontSvc := sellservice.NewStorefrontService(sellRepo, backofficeUserSvc, backofficebrandsvc, productSvc)
	sellBackofficeSvc := sellservice.NewBackofficeService(sellRepo, backofficeUserSvc, backofficebrandsvc, productSvc)
	sellHandler := sellhandler.New(sellBackofficeSvc, sellStorefrontSvc, translator, sellVld, authSvc, cfg.Auth, aclSvc)

	permissionSvc := permissionservice.New(permissionRepo)
	roleSvc := roleservice.New(roleRepo, permissionSvc)
	aclHandler := aclhandler.New(permissionSvc, roleSvc, translator, authSvc, cfg.Auth, aclSvc)

	server := httpserver.New(cfg, backOfficeHandler, brandHandler, mediaHandler, categoryHandler, articleHandler, attrHandler,
		attrDefHandler, productHandler, backofficeHandler, userHandler, commentHandler, discountHandler, cartHandler,
		orderHandler, shippingHandler, reapirHandler, sellHandler, aclHandler)
	server.Serve()
}
