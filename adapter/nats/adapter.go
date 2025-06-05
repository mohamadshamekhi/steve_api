package natsadapter

import (
	"github.com/abolfazlnorzad/steve/pkg/logger"
	"github.com/nats-io/nats.go"
	"time"
)

type Config struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

type Adapter struct {
	nc *nats.Conn            // اتصال به NATS (Core NATS)
	js nats.JetStreamContext // زمینه JetStream
}

func New(cfg Config) *Adapter {

	nc, err := nats.Connect(nats.DefaultURL,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Second*5),
		nats.Timeout(10*time.Second),
		nats.ErrorHandler(func(nc *nats.Conn, s *nats.Subscription, err error) {
			if s != nil {
				logger.L().Error("NATS Async error in %q/%q: %s", s.Subject, s.Queue, err)
			} else {
				logger.L().Error("NATS Async error: %s", err)
			}
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			logger.L().Error("NATS connection closed.")
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.L().Error("NATS reconnected to %s", nc.ConnectedUrl())
		}),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.L().Error("NATS disconnected. Will try to reconnect: %v", err)
		}),
	)
	if err != nil {
		panic("err in connection to nats server :" + err.Error())
	}

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))

	if err != nil {
		nc.Close()
		panic("err in connection to jet stream server :" + err.Error())
	}
	return &Adapter{
		nc: nc,
		js: js,
	}
}
func (a *Adapter) NatsCore() *nats.Conn {
	return a.nc
}
func (a *Adapter) NatsJS() nats.JetStreamContext {
	return a.js
}
