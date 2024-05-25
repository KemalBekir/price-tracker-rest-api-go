package handler

import (
	"encoding/json"
	"net/http"

	"github.com/KemalBekir/price-tracker-rest-api-go/internal/services"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
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

type ScrapeRequest struct {
	URL string `json:"url"`
}

func ScrapeHandler(searchesCollection, pricesCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setJSONHeader(w)

		var scrapeReq ScrapeRequest
		if err := json.NewDecoder(r.Body).Decode(&scrapeReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if scrapeReq.URL == "" {
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		search, err := services.ScrapeAmazon(scrapeReq.URL, searchesCollection, pricesCollection)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(search)
	}
}
