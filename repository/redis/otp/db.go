package redisotp

import redisadapter "github.com/abolfazlnorzad/steve/adapter/redis"

type DB struct {
	adapter *redisadapter.Adapter
}

func New(adapter *redisadapter.Adapter) *DB {
	return &DB{
		adapter: adapter,
	}
}
