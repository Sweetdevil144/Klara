package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/config"
)

var client *mongo.Client = nil
var connected bool = false

func Init() error {
	clientOptions := options.Client().ApplyURI(config.Config("MONGO_URI"))
	var ctx = context.Background()
	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	client = c
	connected = true
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return nil
}

func Connect() (*mongo.Database, error) {
	if !connected {
		return nil, fmt.Errorf("client not connected")
	}

	return client.Database(config.Config("MONGO_DB_NAME")), nil
}

func Disconnect() error {
	if err := client.Disconnect(context.TODO()); err != nil {
		return err
	}
	fmt.Println("connection go brrrr")
	return nil
}
