package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Config struct {
	Host     string `koanf:"host"`
	Port     string `koanf:"port"`
	DBName   string `koanf:"db_name"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}

type DB struct {
	db *mongo.Database
}

func NewMongoDB(cfg Config) *DB {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic("can't connect to mongo db , " + err.Error())
	}

	db := client.Database(cfg.DBName)

	return &DB{
		db: db,
	}
}

func (d *DB) Connection() *mongo.Database {
	return d.db
}
