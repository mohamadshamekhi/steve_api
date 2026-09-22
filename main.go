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
	dashboardhandler "github.com/abolfazlnorzad/steve/delivery/http_server/dashboard"
	discounthandler "github.com/abolfazlnorzad/steve/delivery/http_server/discount"
	inventoryhandler "github.com/abolfazlnorzad/steve/delivery/http_server/inventory"
	mediahandler "github.com/abolfazlnorzad/steve/delivery/http_server/media"
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
	mongodbcash "github.com/abolfazlnorzad/steve/repository/mongo/cash"
	mongodbcategory "github.com/abolfazlnorzad/steve/repository/mongo/category"
	mongodbcomment "github.com/abolfazlnorzad/steve/repository/mongo/comment"
	mongodbcontactus "github.com/abolfazlnorzad/steve/repository/mongo/contactus"
	mongodbdiscount "github.com/abolfazlnorzad/steve/repository/mongo/discount"
	mongodbexpense "github.com/abolfazlnorzad/steve/repository/mongo/expense"
	mongodbinventory "github.com/abolfazlnorzad/steve/repository/mongo/inventory"
	mongodbmedia "github.com/abolfazlnorzad/steve/repository/mongo/media"
	mongodborder "github.com/abolfazlnorzad/steve/repository/mongo/order"
	mongodbpartner "github.com/abolfazlnorzad/steve/repository/mongo/partner"
	mongodbpayment "github.com/abolfazlnorzad/steve/repository/mongo/payment"
	mongodbproduct "github.com/abolfazlnorzad/steve/repository/mongo/product"
	mongodbpurchase "github.com/abolfazlnorzad/steve/repository/mongo/purchase"
	mongodbrating "github.com/abolfazlnorzad/steve/repository/mongo/rating"
	mongodbreferral "github.com/abolfazlnorzad/steve/repository/mongo/referral"
	mongodbrepair "github.com/abolfazlnorzad/steve/repository/mongo/repair"
	mongodbsalesreturn "github.com/abolfazlnorzad/steve/repository/mongo/salesreturn"
	mongodbsell "github.com/abolfazlnorzad/steve/repository/mongo/sell"
	mongodbshipping "github.com/abolfazlnorzad/steve/repository/mongo/shipping"
	mongodbsupplier "github.com/abolfazlnorzad/steve/repository/mongo/supplier"
	mongodbuser "github.com/abolfazlnorzad/steve/repository/mongo/user"
	redisarticleview "github.com/abolfazlnorzad/steve/repository/redis/articleview"
	redisotp "github.com/abolfazlnorzad/steve/repository/redis/otp"
	cartscheduler "github.com/abolfazlnorzad/steve/scheduler/cart"
	orderscheduler "github.com/abolfazlnorzad/steve/scheduler/order"
	accountingservice "github.com/abolfazlnorzad/steve/service/accounting"
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
	dashboardservice "github.com/abolfazlnorzad/steve/service/dashboard"
	discountservice "github.com/abolfazlnorzad/steve/service/discount"
	fetcherservice "github.com/abolfazlnorzad/steve/service/fetcher"
	inventoryservice "github.com/abolfazlnorzad/steve/service/inventory"
	mediaservice "github.com/abolfazlnorzad/steve/service/media"
	orderservice "github.com/abolfazlnorzad/steve/service/order"
	orderdeliveryservice "github.com/abolfazlnorzad/steve/service/orderdelivery"
	ordernotifyservice "github.com/abolfazlnorzad/steve/service/ordernotify"
	partnerservice "github.com/abolfazlnorzad/steve/service/partner"
	paymentservice "github.com/abolfazlnorzad/steve/service/payment"
	productservice "github.com/abolfazlnorzad/steve/service/product"
	purchaseservice "github.com/abolfazlnorzad/steve/service/purchase"
	ratingservice "github.com/abolfazlnorzad/steve/service/rating"
	referralservice "github.com/abolfazlnorzad/steve/service/referral"
	repairservice "github.com/abolfazlnorzad/steve/service/repair"
	salesreturnservice "github.com/abolfazlnorzad/steve/service/salesreturn"
	searchservice "github.com/abolfazlnorzad/steve/service/search"
	sellservice "github.com/abolfazlnorzad/steve/service/sell"
	settingservice "github.com/abolfazlnorzad/steve/service/setting"
	supplierservice "github.com/abolfazlnorzad/steve/service/supplier"
	userservice "github.com/abolfazlnorzad/steve/service/user"
	accountingvalidator "github.com/abolfazlnorzad/steve/validator/accounting"
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
	inventoryvalidator "github.com/abolfazlnorzad/steve/validator/inventory"
	mediavalidator "github.com/abolfazlnorzad/steve/validator/media"
	ordervalidator "github.com/abolfazlnorzad/steve/validator/order"
	orderdeliveryvalidator "github.com/abolfazlnorzad/steve/validator/orderdelivery"
	partnervalidator "github.com/abolfazlnorzad/steve/validator/partner"
	productvalidator "github.com/abolfazlnorzad/steve/validator/product"
	purchasevalidator "github.com/abolfazlnorzad/steve/validator/purchase"
	ratingvalidator "github.com/abolfazlnorzad/steve/validator/rating"
	referralvalidator "github.com/abolfazlnorzad/steve/validator/referral"
	repairvalidator "github.com/abolfazlnorzad/steve/validator/repair"
	salesreturnvalidator "github.com/abolfazlnorzad/steve/validator/salesreturn"
	sellvalidator "github.com/abolfazlnorzad/steve/validator/sell"
	settingvalidator "github.com/abolfazlnorzad/steve/validator/setting"
	suppliervalidator "github.com/abolfazlnorzad/steve/validator/supplier"
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
	permissionRepo := mongodbpermission.New(mongoRepo)
	roleRepo := mongodbrole.New(mongoRepo)
	backOfficeSvc := backofficeservice.New(cfg.BackOffice, authSvc, backOfficeMongoRepo, redisotp, roleRepo)
	backOfficeVld := backofficevalidatorauth.New(backOfficeMongoRepo)

	translator, _ := i18n.New("config/locales", "fa")

	if err := seeder.Run(permissionRepo, roleRepo, backOfficeMongoRepo); err != nil {
		logger.L().Error("seeder failed", "error", err)
	}

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
	referralBackofficeSvc := referralservice.NewBackofficeService(cfg.Referral, referralRepo, userRepo)
	userSvc := userservice.NewStoreFrontService(cfg.User, userRepo, authSvc, redisotp, referralSvc)
	referralHandler := referralhandler.New(referralSvc, referralBackofficeSvc, referralVld, translator, authSvc, cfg.Auth, aclSvc)
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
	// Every status change goes through this decorator so the customer is
	// SMS'd no matter which flow moved the order.
	orderRepo := ordernotifyservice.WrapOrderRepo(mongodborder.New(mongoRepo))
	orderRepo.SetNotifier(ordernotifyservice.New(cfg.OrderNotify, orderRepo, discountRepo, userRepo, mp))
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

	// --- inventory / purchasing / suppliers (انبارداری و حسابداری) ---
	supplierRepo := mongodbsupplier.New(mongoRepo)
	inventoryRepo := mongodbinventory.New(mongoRepo)
	purchaseRepo := mongodbpurchase.New(mongoRepo)

	supplierSvc := supplierservice.NewBackofficeService(supplierRepo)
	inventorySvc := inventoryservice.NewBackofficeService(inventoryRepo, productRepo)
	inventorySvc.SetOrderProvider(mongodborder.New(mongoRepo))
	purchaseSvc := purchaseservice.NewBackofficeService(purchaseRepo, inventorySvc, productRepo, supplierRepo)

	// The order services record every sale in the stock ledger. Wired after
	// construction so the order domain keeps no hard dependency on inventory.
	orderSvc.SetInventoryRecorder(inventorySvc)
	orderBackofficeSvc.SetInventoryRecorder(inventorySvc)

	// --- accounting: expenses, cashbox, returns, financial reports ---
	expenseRepo := mongodbexpense.New(mongoRepo)
	cashRepo := mongodbcash.New(mongoRepo)
	salesReturnRepo := mongodbsalesreturn.New(mongoRepo)

	salesReturnSvc := salesreturnservice.NewBackofficeService(salesReturnRepo, orderRepo, inventorySvc, productRepo)
	accountingSvc := accountingservice.NewBackofficeService(expenseRepo, cashRepo, inventorySvc, salesReturnSvc, purchaseRepo)

	// Everything that moves money tells the cashbox. Wired after
	// construction so no domain depends on accounting at build time.
	salesReturnSvc.SetCashRecorder(accountingSvc)
	purchaseSvc.SetCashRecorder(accountingSvc)
	orderBackofficeSvc.SetCashRecorder(accountingSvc)

	supplierHandler := supplierhandler.New(supplierSvc, translator, suppliervalidator.New(), authSvc, cfg.Auth, aclSvc)
	purchaseHandler := purchasehandler.New(purchaseSvc, translator, purchasevalidator.New(), authSvc, cfg.Auth, aclSvc)
	inventoryHandler := inventoryhandler.New(inventorySvc, translator, inventoryvalidator.New(), authSvc, cfg.Auth, aclSvc)
	accountingHandler := accountinghandler.New(accountingSvc, translator, accountingvalidator.New(), authSvc, cfg.Auth, aclSvc)

	dashboardSvc := dashboardservice.New(accountingSvc, inventoryRepo, cashRepo, productRepo,
		mongodborder.New(mongoRepo), userRepo, repairRepo, aclSvc)
	dashboardHandler := dashboardhandler.New(dashboardSvc, translator, authSvc, cfg.Auth)
	salesReturnHandler := salesreturnhandler.New(salesReturnSvc, translator, salesreturnvalidator.New(), authSvc, cfg.Auth, aclSvc)

	orderDeliverySvc := orderdeliveryservice.New(cfg.OrderDelivery, orderRepo, userRepo, authSvc, redisotp)
	orderDeliveryHandler := orderdeliveryhandler.New(orderDeliverySvc, orderdeliveryvalidator.New(), translator, authSvc, cfg.Auth, aclSvc)

	server := httpserver.New(cfg, backOfficeHandler, brandHandler, mediaHandler, categoryHandler, articleHandler, attrHandler,
		attrDefHandler, productHandler, backofficeHandler, userHandler, commentHandler, discountHandler, cartHandler,
		orderHandler, shippingHandler, reapirHandler, sellHandler, aclHandler, contactHandler, ratingHandler, bannerHandler, paymentHandler, searchHandler, partnerHandler, referralHandler,
		supplierHandler, purchaseHandler, inventoryHandler, accountingHandler, salesReturnHandler, orderDeliveryHandler, dashboardHandler)

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
