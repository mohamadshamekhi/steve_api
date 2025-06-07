package main

import (
	"context"
	"fmt"
	natsadapter "github.com/abolfazlnorzad/steve/adapter/nats"
	redisadapter "github.com/abolfazlnorzad/steve/adapter/redis"
	"github.com/abolfazlnorzad/steve/config"
	backofficeauthparam "github.com/abolfazlnorzad/steve/param/backoffice/auth"
	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
	mongodbbackoffice "github.com/abolfazlnorzad/steve/repository/mongo/backoffice"
	redisotp "github.com/abolfazlnorzad/steve/repository/redis/otp"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	backofficeservice "github.com/abolfazlnorzad/steve/service/backoffice"
)

func main() {
	cfg := config.C()
	repo := redisadapter.New(cfg.Redis)
	nats := natsadapter.New(natsadapter.Config{
		Host: "localhost",
		Port: 4222,
	})
	authSvc := authservice.New(cfg.Auth, redisotp.New(repo), nats)
	mongoRepo := mongodb.NewMongoDB(cfg.Mongo)
	backOfficeMongoRepo := mongodbbackoffice.New(mongoRepo)
	backOfficeSvc := backofficeservice.New(cfg.BackOffice, authSvc, backOfficeMongoRepo)
	resp, err := backOfficeSvc.SendOTPCode(context.Background(), backofficeauthparam.SendLoginCodeRequest{PhoneNumber: "09011216131"})
	fmt.Println("err , resp", err, resp)
}
