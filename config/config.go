package config

import (
	natsadapter "github.com/abolfazlnorzad/steve/adapter/nats"
	redisadapter "github.com/abolfazlnorzad/steve/adapter/redis"
	liarastorage "github.com/abolfazlnorzad/steve/adapter/storage/liara"
	"github.com/abolfazlnorzad/steve/entity"
	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	backofficeservice "github.com/abolfazlnorzad/steve/service/backoffice"
)

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	HTTPServer    HTTPServer               `koanf:"http_server"`
	Auth          authservice.Config       `koanf:"auth"`
	BackOffice    backofficeservice.Config `koanf:"backoffice"`
	Redis         redisadapter.Config      `koanf:"redis"`
	Nats          natsadapter.Config       `koanf:"nats"`
	Mongo         mongodb.Config           `koanf:"mongo"`
	ObjectStorage liarastorage.Config      `koanf:"storage"`
	Media         entity.Config            `koanf:"media"`
}
