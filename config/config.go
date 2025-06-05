package config

import (
	redisadapter "github.com/abolfazlnorzad/steve/adapter/redis"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	backofficeservice "github.com/abolfazlnorzad/steve/service/backoffice"
)

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	HTTPServer HTTPServer               `koanf:"http_server"`
	Auth       authservice.Config       `koanf:"auth"`
	BackOffice backofficeservice.Config `koanf:"backoffice"`
	Redis      redisadapter.Config      `koanf:"redis"`
}
