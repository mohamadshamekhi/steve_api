package main

import (
	"context"
	"fmt"
	natsadapter "github.com/abolfazlnorzad/steve/adapter/nats"
	redisadapter "github.com/abolfazlnorzad/steve/adapter/redis"
	"github.com/abolfazlnorzad/steve/config"
	backofficeauthparam "github.com/abolfazlnorzad/steve/param/backoffice/auth"
	redisotp "github.com/abolfazlnorzad/steve/repository/redis/otp"
	authservice "github.com/abolfazlnorzad/steve/service/auth"
	backofficeservice "github.com/abolfazlnorzad/steve/service/backoffice"
)

func main() {
	cfg := config.C()
	fmt.Printf("%+v\n", cfg)
	repo := redisadapter.New(cfg.Redis)
	nats := natsadapter.New(natsadapter.Config{
		Host: "localhost",
		Port: 4222,
	})
	authSvc := authservice.New(cfg.Auth, redisotp.New(repo), nats)
	backOfficeSvc := backofficeservice.New(cfg.BackOffice, authSvc)
	resp, err := backOfficeSvc.Login(context.Background(), backofficeauthparam.SendLoginCodeRequest{PhoneNumber: "09011216131"})
	fmt.Println("err", err)
	fmt.Println("resp", resp.Message)
}
