package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abolfazlnorzad/steve/adapter/gateway"
	"github.com/abolfazlnorzad/steve/adapter/mellipayamak"
	redisadapter "github.com/abolfazlnorzad/steve/adapter/redis"
	liarastorage "github.com/abolfazlnorzad/steve/adapter/storage/liara"
	"github.com/abolfazlnorzad/steve/adapter/telegrambot"
	"github.com/abolfazlnorzad/steve/config"
	httpserver "github.com/abolfazlnorzad/steve/delivery/http_server"
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
	mediahandler "github.com/abolfazlnorzad/steve/delivery/http_server/media"
	orderhandler "github.com/abolfazlnorzad/steve/delivery/http_server/order"
	partnerhandler "github.com/abolfazlnorzad/steve/delivery/http_server/partner"
	paymenthandler "github.com/abolfazlnorzad/steve/delivery/http_server/payment"
	producthandler "github.com/abolfazlnorzad/steve/delivery/http_server/product"
	ratinghandler "github.com/abolfazlnorzad/steve/delivery/http_server/rating"
	referralhandler "github.com/abolfazlnorzad/steve/delivery/http_server/referral"
	repairhandler "github.com/abolfazlnorzad/steve/delivery/http_server/repair"
	searchhandler "github.com/abolfazlnorzad/steve/delivery/http_server/search"
	sellhandler "github.com/abolfazlnorzad/steve/delivery/http_server/sell"
	shippinghandler "github.com/abolfazlnorzad/steve/delivery/http_server/shipping"
	userhandler "github.com/abolfazlnorzad/steve/delivery/http_server/user"
	"github.com/abolfazlnorzad/steve/internal/seeder"
	"github.com/abolfazlnorzad/steve/pkg/i18n"
	"github.com/abolfazlnorzad/steve/pkg/logger"
	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	mongodbpermission "github.com/abolfazlnorzad/steve/repository/mongo/acl/permission"
	mongodbrole "github.com/abolfazlnorzad/steve/repository/mongo/acl/role"
	mongodbarticle "github.com/abolfazlnorzad/steve/repository/mongo/article"
	mongodbattribute "github.com/abolfazlnorzad/steve/repository/mongo/attribute"
	mongodbattributedefinition "github.com/abolfazlnorzad/steve/repository/mongo/attribute_definition"
	mongodbbackoffice "github.com/abolfazlnorzad/steve/repository/mongo/backoffice"
	mongodbbanner "github.com/abolfazlnorzad/steve/repository/mongo/banner"
	mongodbbrand "github.com/abolfazlnorzad/steve/repository/mongo/brand"
	mongodbcart "github.com/abolfazlnorzad/steve/repository/mongo/cart"
	mongodbcategory "github.com/abolfazlnorzad/steve/repository/mongo/category"
	mongodbcomment "github.com/abolfazlnorzad/steve/repository/mongo/comment"
	mongodbcontactus "github.com/abolfazlnorzad/steve/repository/mongo/contactus"
	mongodbdiscount "github.com/abolfazlnorzad/steve/repository/mongo/discount"
	mongodbmedia "github.com/abolfazlnorzad/steve/repository/mongo/media"
	mongodborder "github.com/abolfazlnorzad/steve/repository/mongo/order"
	mongodbpartner "github.com/abolfazlnorzad/steve/repository/mongo/partner"
	mongodbpayment "github.com/abolfazlnorzad/steve/repository/mongo/payment"
	mongodbproduct "github.com/abolfazlnorzad/steve/repository/mongo/product"
	mongodbrating "github.com/abolfazlnorzad/steve/repository/mongo/rating"
	mongodbreferral "github.com/abolfazlnorzad/steve/repository/mongo/referral"
	mongodbrepair "github.com/abolfazlnorzad/steve/repository/mongo/repair"
	mongodbsell "github.com/abolfazlnorzad/steve/repository/mongo/sell"
	mongodbshipping "github.com/abolfazlnorzad/steve/repository/mongo/shipping"
	mongodbuser "github.com/abolfazlnorzad/steve/repository/mongo/user"
	redisarticleview "github.com/abolfazlnorzad/steve/repository/redis/articleview"
	redisotp "github.com/abolfazlnorzad/steve/repository/redis/otp"
	cartscheduler "github.com/abolfazlnorzad/steve/scheduler/cart"
	orderscheduler "github.com/abolfazlnorzad/steve/scheduler/order"
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
	contactusservice "github.com/abolfazlnorzad/steve/service/contactus"
	discountservice "github.com/abolfazlnorzad/steve/service/discount"
	fetcherservice "github.com/abolfazlnorzad/steve/service/fetcher"
	mediaservice "github.com/abolfazlnorzad/steve/service/media"
	orderservice "github.com/abolfazlnorzad/steve/service/order"
	partnerservice "github.com/abolfazlnorzad/steve/service/partner"
	paymentservice "github.com/abolfazlnorzad/steve/service/payment"
	productservice "github.com/abolfazlnorzad/steve/service/product"
	ratingservice "github.com/abolfazlnorzad/steve/service/rating"
	referralservice "github.com/abolfazlnorzad/steve/service/referral"
	repairservice "github.com/abolfazlnorzad/steve/service/repair"
	searchservice "github.com/abolfazlnorzad/steve/service/search"
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
	contactusvalidator "github.com/abolfazlnorzad/steve/validator/contactus"
	discountvalidator "github.com/abolfazlnorzad/steve/validator/discount"
	mediavalidator "github.com/abolfazlnorzad/steve/validator/media"
	ordervalidator "github.com/abolfazlnorzad/steve/validator/order"
	partnervalidator "github.com/abolfazlnorzad/steve/validator/partner"
	productvalidator "github.com/abolfazlnorzad/steve/validator/product"
	ratingvalidator "github.com/abolfazlnorzad/steve/validator/rating"
	referralvalidator "github.com/abolfazlnorzad/steve/validator/referral"
	repairvalidator "github.com/abolfazlnorzad/steve/validator/repair"
	sellvalidator "github.com/abolfazlnorzad/steve/validator/sell"
	settingvalidator "github.com/abolfazlnorzad/steve/validator/setting"
	uservalidator "github.com/abolfazlnorzad/steve/validator/user"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg := config.C()
	mp := mellipayamak.New(cfg.MelliPayamak)
	adminTelegramBot := telegrambot.NewAdmin(cfg.TelegramBot)
	developerTelegramBot := telegrambot.NewDeveloper(cfg.TelegramBot)
	redisRepo := redisadapter.New(cfg.Redis)
	redisotp := redisotp.New(redisRepo)
	authSvc := authservice.New(cfg.Auth, redisotp, mp)
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

	articleRepo := mongodbarticle.New(mongoRepo)
	articleViewRepo := redisarticleview.New(redisRepo)
	articleVld := articlevalidator.New(articleRepo, articleRepo, mediaSvc)
	articleSvc := articleservice.NewBackofficeService(articleRepo, mediaSvc, backOfficeMongoRepo)
	storefrontSvc := articleservice.NewStorefrontService(articleRepo, mediaSvc, backOfficeMongoRepo, articleRepo, articleViewRepo)
	storefrontVld := articlevalidator.NewStorefrontValidator(articleRepo)
	articleHandler := articlehandler.New(translator, authSvc, cfg.Auth, aclSvc, articleSvc, articleVld, storefrontSvc, storefrontVld)

	attrHandler := attributehandler.New(authSvc, cfg.Auth, aclSvc, translator, attrVld, attrSvc)

	attrDefVld := attributedefinitionvalidator.New(attrDefRepo)
	attrDefSvc := attributeservice.NewAttributeDefinition(attrDefRepo)
	attrDefHandler := attributedefinitionhandler.New(authSvc, cfg.Auth, aclSvc, translator, attrDefVld, attrDefSvc)

	productRepo := mongodbproduct.New(mongoRepo)
	productVld := productvalidator.New(categorySvc, backofficebrandsvc, productRepo, mediaSvc)
	productSvc := productservice.NewBackofficeService(productRepo, mediaSvc, categorySvc, backofficebrandsvc)

	backofficeValidation := backofficevalidator.New(backOfficeMongoRepo)
	backofficeHandler := backofficehandler.New(backOfficeSvc, backofficeValidation, translator, authSvc, cfg.Auth, aclSvc)
	userRepo := mongodbuser.New(mongoRepo)
	userVld := uservalidator.New(userRepo)
	// discountRepo is built here (ahead of the discount block below) so the
	// referral service can mint reward/welcome coupons, and referralSvc is
	// wired into the storefront user service for signup-time attribution.
	discountRepo := mongodbdiscount.New(mongoRepo)
	referralRepo := mongodbreferral.New(mongoRepo)
	referralVld := referralvalidator.New()
	referralSvc := referralservice.NewStorefrontService(cfg.Referral, referralRepo, discountRepo, userRepo, mp)
	userSvc := userservice.NewStoreFrontService(cfg.User, userRepo, authSvc, redisotp, referralSvc)
	referralHandler := referralhandler.New(referralSvc, referralVld, translator, authSvc, cfg.Auth)
	backofficeUserSvc := userservice.NewBackofficeService(userRepo)
	userHandler := userhandler.New(backofficeUserSvc, userSvc, userVld, translator, authSvc, cfg.Auth, aclSvc)
	commentRepo := mongodbcomment.New(mongoRepo)
	commentVld := commentvalidator.New(commentRepo)

	commentableFetecher := fetcherservice.NewCommentableFetcher(*productRepo, *articleRepo)
	commenterFetecher := fetcherservice.NewCommenterFetcher(*userRepo, *backOfficeMongoRepo)
	commentSvc := commentservice.New(commentRepo, commenterFetecher, commentableFetecher)
	storefrontCommentSvc := commentservice.NewStorefrontService(commentRepo, commentableFetecher, commenterFetecher)
	commentHandler := commenthandler.New(commentSvc, translator, commentVld, authSvc, cfg.Auth, aclSvc, storefrontCommentSvc)
	discountVld := discountvalidator.New(discountRepo, categoryRepo, productRepo, userRepo)

	cartRepo := mongodbcart.New(mongoRepo)
	cartVld := cartvalidator.New(productRepo, cartRepo, userRepo)
	cartSvc := cartservice.NewStorefrontService(cartRepo, productRepo, userRepo)
	discountbackofficeSvc := discountservice.NewBackofficeService(discountRepo, &cartSvc, productSvc)
	productStorefrontSvc := productservice.NewStorefrontService(productRepo, mediaSvc, categorySvc, backofficebrandsvc, discountbackofficeSvc)
	productHandler := producthandler.NewHandler(translator, authSvc, cfg.Auth, aclSvc, productSvc, productStorefrontSvc, productVld)

	discountHandler := discounthandler.New(discountbackofficeSvc, translator, discountVld, authSvc, cfg.Auth, aclSvc)
	cartBackofficeSvc := cartservice.NewBackofficeService(cartRepo, userRepo, productRepo, productSvc)
	orderRepo := mongodborder.New(mongoRepo)
	orderVld := ordervalidator.New(discountbackofficeSvc, cartRepo)
	shippingRepo := mongodbshipping.New(mongoRepo)
	shippingVld := settingvalidator.NewShippingValidator(shippingRepo)
	shippingBakcofficeSvc := settingservice.NewBackofficeShippingService(shippingRepo)
	shippingStorefrontSvc := settingservice.NewStorefrontShippingService(shippingRepo)
	orderSvc := orderservice.New(&cartSvc, discountbackofficeSvc, productSvc, orderRepo, userSvc, shippingBakcofficeSvc)
	cartSvc.SetOrderProvider(orderSvc)
	cartHandler := carthandler.New(cartSvc, cartBackofficeSvc, orderSvc, cartVld, translator, authSvc, cfg.Auth, aclSvc)
	storefrontCategorySvc := categoryservice.NewStorefrontService(categoryRepo, mediaSvc, attrSvc, productStorefrontSvc)
	categoryHandler := categoryhandler.New(categorySvc, categoryVld, translator, authSvc, cfg.Auth, aclSvc, storefrontCategorySvc)
	orderBackofficeSvc := orderservice.NewBackofficeService(orderRepo, backofficeUserSvc, productSvc, discountbackofficeSvc, &cartBackofficeSvc)

	shippingHandler := shippinghandler.New(shippingBakcofficeSvc, shippingStorefrontSvc, shippingVld, translator, authSvc, cfg.Auth, aclSvc)

	repairRepo := mongodbrepair.New(mongoRepo)
	repairVld := repairvalidator.New()
	repairStorfrontSvc := repairservice.NewStorefrontService(repairRepo, backofficeUserSvc, backofficebrandsvc, productSvc, backOfficeSvc, mp)
	repairBackofficeSvc := repairservice.NewBackofficeService(repairRepo, backofficeUserSvc, backofficebrandsvc, productSvc)
	reapirHandler := repairhandler.New(repairBackofficeSvc, repairStorfrontSvc, translator, repairVld, authSvc, cfg.Auth, aclSvc)

	sellRepo := mongodbsell.New(mongoRepo)
	sellVld := sellvalidator.New()
	sellStorefrontSvc := sellservice.NewStorefrontService(sellRepo, backofficeUserSvc, backofficebrandsvc, productSvc, backOfficeSvc, mp)
	sellBackofficeSvc := sellservice.NewBackofficeService(sellRepo, backofficeUserSvc, backofficebrandsvc, productSvc)
	sellHandler := sellhandler.New(sellBackofficeSvc, sellStorefrontSvc, translator, sellVld, authSvc, cfg.Auth, aclSvc)

	permissionSvc := permissionservice.New(permissionRepo)
	roleSvc := roleservice.New(roleRepo, permissionSvc)
	aclHandler := aclhandler.New(permissionSvc, roleSvc, translator, authSvc, cfg.Auth, aclSvc)

	contactRepo := mongodbcontactus.New(mongoRepo)
	contactVld := contactusvalidator.New()
	storefrontContactUsSvc := contactusservice.NewStorefrontService(contactRepo)
	backofficeContactUsSvc := contactusservice.NewBackofficeService(contactRepo)
	contactHandler := contactushandler.New(storefrontContactUsSvc, translator, contactVld, authSvc, cfg.Auth, aclSvc, backofficeContactUsSvc)
	ratingRepo := mongodbrating.New(mongoRepo)
	ratingVld := ratingvalidator.New()
	ratingSvc := ratingservice.New(ratingRepo, productStorefrontSvc)
	ratingHandler := ratinghandler.New(ratingSvc, translator, ratingVld, authSvc, cfg.Auth, aclSvc)
	bannerRepo := mongodbbanner.New(mongoRepo)
	bannerBackofficeSvc := settingservice.NewBannerBackofficeService(bannerRepo, mediaSvc)
	bannerVld := settingvalidator.NewBannerValidator()
	bannerHandler := bannerhandler.New(bannerBackofficeSvc, bannerVld, translator, authSvc, cfg.Auth, aclSvc)
	paymentRepo := mongodbpayment.New(mongoRepo)
	zarinpal := gateway.NewZarinpalAdapter(cfg.Zarinpal)
	paymentSvc := paymentservice.NewService(cfg.Payment, paymentRepo, orderSvc, zarinpal, discountbackofficeSvc, productStorefrontSvc, mp, userSvc, adminTelegramBot, developerTelegramBot, backOfficeSvc)
	paymentHandler := paymenthandler.New(paymentSvc, translator, authSvc, cfg.Auth, aclSvc)
	orderSvc.SetPaymentProvider(paymentSvc)
	orderBackofficeSvc.SetPaymentProvider(paymentSvc)
	orderHandler := orderhandler.New(orderSvc, orderBackofficeSvc, orderVld, translator, authSvc, cfg.Auth, aclSvc)

	searchSvc := searchservice.NewStorefrontService(productStorefrontSvc, storefrontSvc)
	searchHandler := searchhandler.New(translator, searchSvc)

	partnerRepo := mongodbpartner.New(mongoRepo)
	partnerResolver := partnerservice.NewPartnerPriceResolver()
	partnerStorefrontSvc := partnerservice.NewStorefrontService(partnerRepo, cfg.Partner, adminTelegramBot)
	partnerShopSvc := partnerservice.NewShopService(partnerRepo, productRepo, orderRepo, cartRepo, orderRepo, productSvc, paymentSvc, partnerResolver, cfg.Partner)
	partnerBackofficeSvc := partnerservice.NewBackofficeService(partnerRepo, productRepo, userRepo, backOfficeMongoRepo, orderRepo, paymentSvc, mp, cfg.Partner)
	partnerVld := partnervalidator.New()
	partnerHandler := partnerhandler.New(partnerStorefrontSvc, partnerShopSvc, partnerBackofficeSvc, partnerVld, translator, authSvc, cfg.Auth, aclSvc)

	server := httpserver.New(cfg, backOfficeHandler, brandHandler, mediaHandler, categoryHandler, articleHandler, attrHandler,
		attrDefHandler, productHandler, backofficeHandler, userHandler, commentHandler, discountHandler, cartHandler,
		orderHandler, shippingHandler, reapirHandler, sellHandler, aclHandler, contactHandler, ratingHandler, bannerHandler, paymentHandler, searchHandler, partnerHandler, referralHandler)

	orderScheduler := orderscheduler.New(cfg.OrderScheduler, orderBackofficeSvc)
	cartScheduler := cartscheduler.New(cfg.CartScheduler, cartBackofficeSvc)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		orderScheduler.OrderSchedulerRun(ctx)
		return nil
	})

	g.Go(func() error {
		cartScheduler.CartSchedulerRun(ctx)
		return nil
	})

	g.Go(func() error {
		server.Serve()
		return nil
	})

	<-ctx.Done()
	fmt.Println("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	server.Shutdown(shutdownCtx)

	if err := g.Wait(); err != nil {
		logger.L().Error("application stopped with error", "error", err)
	}
}
