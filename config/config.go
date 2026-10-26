package config

import (
	"github.com/abolfazlnorzad/steve/adapter/gateway"
	"github.com/abolfazlnorzad/steve/adapter/mellipayamak"
	redisadapter "github.com/abolfazlnorzad/steve/adapter/redis"
	liarastorage "github.com/abolfazlnorzad/steve/adapter/storage/liara"
	"github.com/abolfazlnorzad/steve/adapter/telegrambot"
	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	birthdayscheduler "github.com/abolfazlnorzad/steve/scheduler/birthday"
	cartscheduler "github.com/abolfazlnorzad/steve/scheduler/cart"
	"github.com/abolfazlnorzad/steve/scheduler/order"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	backofficeservice "github.com/abolfazlnorzad/steve/service/backoffice"
	birthdayservice "github.com/abolfazlnorzad/steve/service/birthday"
	mediaservice "github.com/abolfazlnorzad/steve/service/media"
	orderdeliveryservice "github.com/abolfazlnorzad/steve/service/orderdelivery"
	ordernotifyservice "github.com/abolfazlnorzad/steve/service/ordernotify"
	partnerservice "github.com/abolfazlnorzad/steve/service/partner"
	paymentservice "github.com/abolfazlnorzad/steve/service/payment"
	referralservice "github.com/abolfazlnorzad/steve/service/referral"
	userservice "github.com/abolfazlnorzad/steve/service/user"
)

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	HTTPServer        HTTPServer                  `koanf:"http_server"`
	Auth              authservice.Config          `koanf:"auth"`
	BackOffice        backofficeservice.Config    `koanf:"backoffice"`
	User              userservice.Config          `koanf:"user"`
	Redis             redisadapter.Config         `koanf:"redis"`
	Mongo             mongodb.Config              `koanf:"mongo"`
	ObjectStorage     liarastorage.Config         `koanf:"object_storage"`
	Media             mediaservice.Config         `koanf:"media"`
	MelliPayamak      mellipayamak.Config         `koanf:"mellipayamak"`
	Payment           paymentservice.Config       `koanf:"payment"`
	Zarinpal          gateway.Config              `koanf:"zarinpal"`
	OrderScheduler    orderscheduler.Config       `koanf:"order_scheduler"`
	CartScheduler     cartscheduler.Config        `koanf:"cart_scheduler"`
	TelegramBot       telegrambot.Config          `koanf:"telegram_bot"`
	Partner           partnerservice.Config       `koanf:"partner"`
	Referral          referralservice.Config      `koanf:"referral"`
	OrderNotify       ordernotifyservice.Config   `koanf:"order_notify"`
	OrderDelivery     orderdeliveryservice.Config `koanf:"order_delivery"`
	Birthday          birthdayservice.Config      `koanf:"birthday"`
	BirthdayScheduler birthdayscheduler.Config    `koanf:"birthday_scheduler"`
}
