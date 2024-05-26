package services

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/KemalBekir/price-tracker-rest-api-go/internal/db"
	"github.com/KemalBekir/price-tracker-rest-api-go/internal/model"
	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO - fix not fetching data correctly
func GetAll(searchesCollection, pricesCollection *mongo.Collection) ([]model.SearchResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Increased timeout
	defer cancel()

	cursor, err := searchesCollection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error finding documents: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var searchResponses []model.SearchResponse
	for cursor.Next(ctx) {
		var search model.Searches
		if err := cursor.Decode(&search); err != nil {
			log.Printf("Error decoding search document: %v", err)
			return nil, err
		}

		var prices []model.Price
		for _, priceID := range search.Prices {
			var price model.Price
			err := pricesCollection.FindOne(ctx, bson.M{"_id": priceID}).Decode(&price)
			if err != nil {
				log.Printf("Error finding price document: %v", err)
				continue
			}
			prices = append(prices, price)
		}

		searchResponse := model.SearchResponse{
			ID:        search.ID,
			URL:       search.URL,
			Domain:    search.Domain,
			ItemName:  search.ItemName,
			Img:       search.Img,
			Prices:    prices,
			Owner:     search.Owner,
			CreatedAt: search.CreatedAt,
			UpdatedAt: search.UpdatedAt,
		}

		searchResponses = append(searchResponses, searchResponse)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, err
	}

	return searchResponses, nil
}

func GetItemByID(id string) (model.Searches, error) {
	collection := db.GetCollection("searches")
	ctx := context.TODO()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Searches{}, err
	}

	var item model.Searches
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&item)
	if err != nil {
		return model.Searches{}, err
	}

	return item, nil
}

func DeleteItemById(id string) error {
	collection := db.GetCollection("searches")
	ctx := context.TODO()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": objID})

	return err
}

func GetAllByOwner(ownerID primitive.ObjectID) ([]model.Searches, error) {
	collection := db.GetCollection("searches")
	ctx := context.TODO()

	cursor, err := collection.Find(ctx, bson.M{"owner": ownerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var items []model.Searches
	if err = cursor.All(ctx, &items); err != nil {
		return nil, err
	}

	return items, nil

}

func CreateSearch(search model.Searches) (model.Searches, error) {
	collection := db.GetCollection("searches")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	search.ID = primitive.NewObjectID()
	search.CreatedAt = time.Now()
	search.UpdatedAt = time.Now()

	_, err := collection.InsertOne(ctx, search)
	return search, err
}

func GetSearchByID(id string) (model.Searches, error) {
	collection := db.GetCollection("searches")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Searches{}, err
	}

	var search model.Searches
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&search)

	return search, err
}

// TODO: fix not getting productName and price
func ScrapeAmazon(url string, searchesCollection, pricesCollection *mongo.Collection) (*model.Searches, error) {
	c := colly.NewCollector()

	var productName, imgSrc string
	var formattedPrice float64

	c.Wait()

	c.OnHTML(".product-hero__title", func(e *colly.HTMLElement) {
		productName = strings.TrimSpace(e.Text)
	})

	c.OnHTML(".inc-vat .price", func(e *colly.HTMLElement) {
		// Extract the price text content
		priceText := e.Text
		// Clean up the price text
		priceText = strings.TrimSpace(priceText)
		priceText = strings.ReplaceAll(priceText, "£", "")
		priceText = strings.ReplaceAll(priceText, ",", "")
		priceText = strings.ReplaceAll(priceText, "inc.", "")
		priceText = strings.ReplaceAll(priceText, "vat", "")
		priceText = strings.ReplaceAll(priceText, " ", "")
		priceText = strings.ReplaceAll(priceText, "\n", "")

		// Convert the cleaned price text to float64
		price, err := strconv.ParseFloat(priceText, 64)
		if err == nil {
			formattedPrice = price
		} else {
			fmt.Println("Error parsing price:", err)
		}
	})

	c.OnHTML(".js-carousel-img", func(e *colly.HTMLElement) {
		imgSrc = e.Attr("src")
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	domain := "ebuyer.com" // Assuming the domain is always ebuyer.com

	search := &model.Searches{
		URL:      url,
		Domain:   domain,
		ItemName: productName,
		Img:      imgSrc,
	}

	// Define a new Price object
	newPrice := model.Price{
		ID:        primitive.NewObjectID(),
		Price:     formattedPrice,
		CreatedAt: time.Now(), // Set the creation time of the price
	}

	filter := bson.M{"url": url}
	update := bson.M{
		"$set": search,
		"$addToSet": bson.M{
			"prices": newPrice,
		},
		"$currentDate": bson.M{
			"updatedAt": true,
		},
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var updatedSearch model.Searches
	err = searchesCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedSearch)
	if err != nil {
		return nil, fmt.Errorf("could not insert or update search: %v", err)
	}

	price := &model.Price{
		Price:     formattedPrice,
		SearchID:  updatedSearch.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = pricesCollection.InsertOne(ctx, price)
	if err != nil {
		return nil, fmt.Errorf("could not insert price: %v", err)
	}

	return &updatedSearch, nil
}
