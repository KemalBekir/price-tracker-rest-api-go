package main

import (
	"log"
	"net/http"

	"github.com/KemalBekir/price-tracker-rest-api-go/internal/db"
	"github.com/KemalBekir/price-tracker-rest-api-go/internal/router"
)

func main() {
	client, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	searchesCollection := client.Database("priceTracker").Collection("searches")
	pricesCollection := client.Database("priceTracker").Collection("pricehistories")

	r := router.SetupRouter(searchesCollection, pricesCollection)
	log.Println("Starting server on :5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}
