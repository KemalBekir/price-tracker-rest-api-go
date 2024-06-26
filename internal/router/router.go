package router

import (
	"github.com/KemalBekir/price-tracker-rest-api-go/internal/handler"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(searchesCollection, pricesCollection *mongo.Collection) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.GetAllItemsHandler).Methods("GET")
	r.HandleFunc("/scrape", handler.ScrapeHandler(searchesCollection, pricesCollection)).Methods("POST")
	r.HandleFunc("/{id}", handler.GetItemHandler).Methods("GET")
	// r.HandleFunc("/{id}", middleware.IsOwner(handler.DeleteItemHandler)).Methods("DELETE")
	// r.HandleFunc("/mySearches", middleware.isAuth(handler.GetUserSearchesHandler)).Methods("GET")
	return r
}
