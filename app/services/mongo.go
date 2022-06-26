package services

import (
	"context"
	"os"
	"time"

	"github.com/google/martian/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectMongo(timeout time.Duration) *mongo.Database {
	uri := os.Getenv("MONGO_URI")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	clientOptions := options.Client()
	clientOptions.ApplyURI(uri)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Error(err)
		panic(err)
	}

	return client.Database(os.Getenv("MONGO_DB_NAME"))
}
