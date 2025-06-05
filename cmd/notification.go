package main

import (
	"fmt"
	natsadapter "github.com/abolfazlnorzad/steve/adapter/nats"
	"github.com/abolfazlnorzad/steve/contract/goproto/auth"
	"github.com/abolfazlnorzad/steve/pkg/constant"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

func main() {

	a := natsadapter.New(natsadapter.Config{
		Host: "localhost",
		Port: 4222,
	})

	_, e := a.NatsJS().Subscribe(constant.SendOTPSubject, func(msg *nats.Msg) {
		fmt.Printf("Received OTP subject: %s\n", msg.Subject)
		var aaa auth.SendOTPCode
		e := proto.Unmarshal(msg.Data, &aaa)
		if e != nil {
			fmt.Println("eeeeee", e)
		}
		fmt.Println(aaa.Code)
	})
	if e != nil {
		fmt.Println("é", e)
	}
}
