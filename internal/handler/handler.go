package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KemalBekir/price-tracker-rest-api-go/internal/services"
	"github.com/gorilla/mux"
)

func setJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func GetAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	setJSONHeader(w)
	data, err := services.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(data)
}

func GetItemHandler(w http.ResponseWriter, r *http.Request) {
	setJSONHeader(w)
	id := mux.Vars(r)["id"]
	item, err := services.GetItemByID(id)
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(item)
}

// func ScrapeHandler(w http.ResponseWriter, r *http.Request) {
// 	setJSONHeader(w)
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
