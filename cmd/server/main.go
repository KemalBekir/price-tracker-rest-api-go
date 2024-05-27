package main

import (
	"log"
	"net/http"

	"github.com/KemalBekir/price-tracker-rest-api-go/internal/db"
	"github.com/KemalBekir/price-tracker-rest-api-go/internal/router"
	"github.com/rs/cors"
)

func main() {
	client, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	searchesCollection := client.Database("priceTracker").Collection("searches")
	pricesCollection := client.Database("priceTracker").Collection("pricehistories")

	r := router.SetupRouter(searchesCollection, pricesCollection)

	// Simplified CORS for debugging
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://pricetracker-api.onrender.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)
	log.Println("Starting server on :5000")
	log.Fatal(http.ListenAndServe(":5000", handler))
}
