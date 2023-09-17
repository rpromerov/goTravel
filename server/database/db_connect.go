package db_connector

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Close(client mongo.Client, ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	defer client.Disconnect(ctx)
}
func Connect() (mongo.Client, context.Context, context.CancelFunc, error) {
	uri := os.Getenv("CONNECTION_STRING")
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return *client, ctx, cancel, err

}
func Ping(client mongo.Client, ctx context.Context) error {
	fmt.Println("ping")

	err := client.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Error pinging MongoDB " + err.Error())
	}
	fmt.Println("pong")

	return err
}
func InsertOne(client mongo.Client, ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	coll := client.Database("gotravel").Collection(collection)
	result, err := coll.InsertOne(ctx, document)
	return result, err
}
func GetOne(client mongo.Client, ctx context.Context, collection string, filter interface{}) (mongo.SingleResult, error) {
	coll := client.Database("gotravel").Collection(collection)
	result := coll.FindOne(ctx, filter)
	return *result, nil
}
