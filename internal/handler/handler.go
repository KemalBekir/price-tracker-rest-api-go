package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/KemalBekir/price-tracker-rest-api-go/internal/db"
	"github.com/KemalBekir/price-tracker-rest-api-go/internal/model"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	collection := db.GetCollection("searches")
	ctx := context.TODO()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Failed to fetch items %s", http.StatusInternalServerError)
	}

	defer cursor.Close(ctx)

	var items []model.Searches
	if err = cursor.All(ctx, &items); err != nil {
		log.Printf("Failed to decode items %s", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(items)
}

// func ScrapeHandler(w http.ResponseWriter, r *http.Request) {
// 	var req struct {
// 		URL    string `json:"url"`
// 		Domain string `json:"domain"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}

// 	price, err := scraper.ScrapePrice(req.URL)
// 	if err != nil {
// 		http.Error(w, "Failed to scrape the price", http.StatusInternalServerError)
// 		return
// 	}
// 	item := model.Searches{URL: req.URL, PRICES: price}
// 	collection := db.GetCollection("searches")

// 	res, err := collection.InsertOne(context.TODO(), item)
// 	if err != nil {
// 		http.Error(w, "Failed to save item to database", http.StatusInternalServerError)
// 		return
// 	}

// 	item.ID = res.InsertedID.(primitive.ObjectID)
// 	json.NewEncoder(w).Encode(item)
// }

// func GetItemHandler(w http.ResponseWriter, r *http.Request) {
// 	id := mux.Vars(r)["id"]
// 	item, err := services.GetItemByID(id)
// 	if err != nil {
// 		http.Error(w, "Item not found", http.StatusNotFound)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(item)
// }
