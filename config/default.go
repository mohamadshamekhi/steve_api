package config

import (
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	partnerservice "github.com/abolfazlnorzad/steve/service/partner"
	referralservice "github.com/abolfazlnorzad/steve/service/referral"
)

func Default() Config {
	return Config{
		HTTPServer: HTTPServer{
			Port: 7777,
		},
		Auth: authservice.Config{
			AccessExpirationTime: BackOfficeAccessTokenExpireDuration,
			AccessSubject:        AccessTokenSubject,
		},
		Partner: partnerservice.Config{
			Enabled:               true,
			AllowDiscountStacking: false,
			AllProductsEligible:   true,
		},
		Referral: referralservice.Config{
			Enabled:                true,
			InviteBaseURL:          "https://stevemobile.ir/invite",
			RewardThreshold:        5,
			RewardDiscountType:     "percentage",
			RewardDiscountValue:    10,
			RewardValidDays:        30,
			WelcomeDiscountEnabled: true,
			WelcomeDiscountType:    "percentage",
			WelcomeDiscountValue:   5,
			WelcomeValidDays:       30,
		},
	}
}
