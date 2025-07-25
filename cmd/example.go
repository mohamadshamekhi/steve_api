package cmd

import (
	"context"
	"fmt"
	"github.com/abolfazlnorzad/steve/config"
	"github.com/abolfazlnorzad/steve/entity"
	mongodb "github.com/abolfazlnorzad/steve/repository/mongo"
)

func main() {
	cfg := config.C()

	mongoRepo := mongodb.NewMongoDB(cfg.Mongo)
	res, err := mongoRepo.Connection().Collection("backofficeUsers").InsertOne(context.Background(), entity.BackOffice{
		ID:           "01G65Z755AFWAKHE12NY0CQ9FH",
		Name:         "nrz",
		PhoneNumber:  "09011216131",
		NationalCode: "7777777",
		Address:      "amol",
		ZipCode:      "1234567",
		AvatarID:     0,
	})
	fmt.Println("err ", err)
	fmt.Println("res ", res)
}
