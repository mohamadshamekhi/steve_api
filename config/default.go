package config

import (
	"time"

	authservice "github.com/abolfazlnorzad/steve/service/auth"
	orderdeliveryservice "github.com/abolfazlnorzad/steve/service/orderdelivery"
	ordernotifyservice "github.com/abolfazlnorzad/steve/service/ordernotify"
	partnerservice "github.com/abolfazlnorzad/steve/service/partner"
	referralservice "github.com/abolfazlnorzad/steve/service/referral"
)

func Default() Config {
	return Config{
		HTTPServer: HTTPServer{
			Port: 7777,
		},
		Auth: authservice.Config{
			AccessExpirationTime:       BackOfficeAccessTokenExpireDuration,
			AccessSubject:              AccessTokenSubject,
			ServiceTokenExpirationTime: ServiceTokenExpireDuration,
		},
		Partner: partnerservice.Config{
			Enabled:               true,
			AllowDiscountStacking: false,
			AllProductsEligible:   true,
		},
		Referral: referralservice.Config{
			Enabled:                true,
			InviteBaseURL:          "https://stevemobile.com/invite",
			RewardThreshold:        5,
			RewardDiscountType:     "percentage",
			RewardDiscountValue:    10,
			RewardValidDays:        30,
			WelcomeDiscountEnabled: true,
			WelcomeDiscountType:    "percentage",
			WelcomeDiscountValue:   5,
			WelcomeValidDays:       30,
		},
		OrderNotify: ordernotifyservice.Config{
			Enabled:                 true,
			ThankYouDiscountEnabled: true,
			ThankYouDiscountType:    "fixed_amount",
			ThankYouDiscountValue:   100000,
			ThankYouValidDays:       30,
			SMSTimeout:              30 * time.Second,
		},
		OrderDelivery: orderdeliveryservice.Config{
			BaseURL:         "https://stevemobile.com/order/delivery",
			OTPPrefix:       "order_delivery",
			MaxOTPAttempts:  5,
			AttemptWindow:   15 * time.Minute,
			QRSize:          512,
			InvoiceTokenTTL: 30 * time.Minute,
		},
	}
}
