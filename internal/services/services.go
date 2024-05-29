package services

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/KemalBekir/price-tracker-rest-api-go/internal/db"
	"github.com/KemalBekir/price-tracker-rest-api-go/internal/model"
	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAll() ([]model.Searches, error) {
	collection := db.GetCollection("searches")
	ctx := context.TODO()

	cursor, err := collection.Find(ctx, bson.M{})
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

func Scrape(url string, searchesCollection, pricesCollection *mongo.Collection) (*model.Searches, error) {
	c := colly.NewCollector()

	var productName, imgSrc string
	var formattedPrice float64

	c.Wait()

	c.OnHTML(".product-hero__title", func(e *colly.HTMLElement) {
		productName = strings.TrimSpace(e.Text)
	})

	c.OnHTML("div.purchase-info__price div.inc-vat p.price", func(e *colly.HTMLElement) {
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

	c.OnHTML("div.image-gallery__hero img", func(e *colly.HTMLElement) {
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

func fetchPrice(url string) (float64, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("ebuyer.com"),
	)

	var price float64

	c.OnHTML("div.purchase-info__price div.inc-vat p.price", func(e *colly.HTMLElement) {
		priceText := e.Text

		priceText = strings.TrimSpace(priceText)
		priceText = strings.ReplaceAll(priceText, "£", "")
		priceText = strings.ReplaceAll(priceText, ",", "")
		priceText = strings.ReplaceAll(priceText, "inc.", "")
		priceText = strings.ReplaceAll(priceText, "vat", "")
		priceText = strings.ReplaceAll(priceText, " ", "")
		priceText = strings.ReplaceAll(priceText, "\n", "")

		var err error
		price, err = strconv.ParseFloat(priceText, 64)
		if err != nil {
			log.Println("Error parsing price: ", err)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL: ", r.Request.URL, "failed with response: ", r, "\nError", err)
	})
	err := c.Visit(url)
	if err != nil {
		return 0, err
	}

	c.Wait()
	return price, nil
}

func UpdatePrices(searchesCollection, pricesCollection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	cursor, err := searchesCollection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("could not fetch searches: %v", err)
	}
	defer cursor.Close(ctx)

	var wg sync.WaitGroup

	for cursor.Next(ctx) {
		var search model.Searches
		if err := cursor.Decode(&search); err != nil {
			return fmt.Errorf("could not decode search: %v", err)
		}

		wg.Add(1)
		go func(search model.Searches) {
			defer wg.Done()

			price, err := fetchPrice(search.URL)
			if err != nil {
				log.Printf("Error fetching price for URL %s: %v", search.URL, err)
				return
			}

			newPrice := model.Price{
				ID:        primitive.NewObjectID(),
				Price:     price,
				CreatedAt: time.Now(),
			}

			filter := bson.M{"_id": search.ID}
			update := bson.M{
				"$addToSet": bson.M{
					"prices": newPrice,
				},
				"$currentDate": bson.M{
					"updatedAt": true,
				},
			}

			opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

			err = searchesCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&search)
			if err != nil {
				log.Printf("Could not update search %s: %v", search.URL, err)
			}

			_, err = pricesCollection.InsertOne(ctx, newPrice)
			if err != nil {
				log.Printf("Could not insert new price for URL %s: %v", search.URL, err)
			}
		}(search)
	}

	wg.Wait()
	return nil
}
