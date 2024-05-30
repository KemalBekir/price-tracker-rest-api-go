package main

import (
	"log"
	"net/http"

	"github.com/KemalBekir/price-tracker-rest-api-go/internal/db"
	"github.com/KemalBekir/price-tracker-rest-api-go/internal/router"
	"github.com/KemalBekir/price-tracker-rest-api-go/internal/services"
	"github.com/robfig/cron/v3"
	"github.com/rs/cors"
)

func main() {
	client, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	searchesCollection := client.Database("priceTracker").Collection("searches")
	pricesCollection := client.Database("priceTracker").Collection("pricehistories")

	c := cron.New()
	_, err = c.AddFunc("0 23 * * *", func() {
		err := services.UpdatePrices(searchesCollection, pricesCollection)
		if err != nil {
			log.Printf("Error updating pricing: %v", err)
		} else {
			log.Println("Price update completed successfully.")
		}
	})
	if err != nil {
		log.Fatalf("Error shceduling daily price update: %v", err)
	}
	c.Start()

	defer c.Stop()

	r := router.SetupRouter(searchesCollection, pricesCollection)
	// Simplified CORS for debugging
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://pricetracker-api.onrender.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(r)
	log.Println("Starting server on :5000")
	log.Fatal(http.ListenAndServe(":5000", handler))
}
