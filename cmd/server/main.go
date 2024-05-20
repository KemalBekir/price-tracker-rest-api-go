package main

import (
	"log"
	"net/http"

	"github.com/KemalBekir/price-tracker-rest-api-go/internal/db"
	"github.com/KemalBekir/price-tracker-rest-api-go/internal/router"
)

func main() {
	db.Connect()

	r := router.SetupRouter()
	log.Fatal(http.ListenAndServe(":5000", r))
}
