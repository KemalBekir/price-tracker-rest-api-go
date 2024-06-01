package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func ConnectDB() (*mongo.Client, error) {
	// err := godotenv.Load("../../.env")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading .env file!")
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
	fmt.Println("Connected to your DB")
	Client = client

	return Client, nil
}

func GetCollection(collectionName string) *mongo.Collection {
	return Client.Database("priceTracker").Collection(collectionName)
}

type ErrorResponse struct {
	StatusCode   int    `json:"status"`
	ErrorMessage string `json:"message"`
}

func GetError(err error, w http.ResponseWriter) {
	log.Fatal(err.Error())
	var response = ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   http.StatusInternalServerError,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}
