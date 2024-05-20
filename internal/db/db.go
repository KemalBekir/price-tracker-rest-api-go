package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading .env file")
	}

	uri := os.Getenv("DB_URI")
	if uri == "" {
		log.Fatal("URI not set in the .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("Error while connecting to your DB")
	}

	Client = client
}

func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("priceTracker").Collection(collectionName)
}
