package services

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/KemalBekir/price-tracker-rest-api-go/internal/db"
	"github.com/KemalBekir/price-tracker-rest-api-go/internal/model"
	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	SearchesCollection *mongo.Collection
}

func NewRepository(client *mongo.Client) *Repository {
	return &Repository{
		SearchesCollection: client.Database("priceTracker").Collection("searches"),
	}
}

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

func ScrapeAndSave(url string, repo *Repository) (*model.Searches, error) {
	c := colly.NewCollector(
		colly.AllowedDomains("amazon.co.uk"),
	)

	var productName, priceStr, imgSrc string

	c.OnHTML("#title", func(e *colly.HTMLElement) {
		productName = strings.TrimSpace(e.Text)
	})

	c.OnHTML(".a-price-whole", func(e *colly.HTMLElement) {
		priceStr = strings.TrimSpace(e.Text)
	})

	c.OnHTML(".a-price-fraction", func(e *colly.HTMLElement) {
		priceStr += "." + strings.TrimSpace(e.Text)
	})

	c.OnHTML(".a-dynamic-image", func(e *colly.HTMLElement) {
		imgSrc = e.Attr("src")
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	if productName == "" || priceStr == "" {
		return nil, errors.New("could not extract necessary information")
	}

	price, err := strconv.ParseFloat(strings.Replace(priceStr, ",", "", -1), 64)
	if err != nil {
		return nil, err
	}

	search := model.Searches{}
	err = repo.SearchesCollection.FindOne(context.TODO(), bson.M{"url": url}).Decode(&search)
	if err == mongo.ErrNoDocuments {
		search = model.Searches{
			ID:        primitive.NewObjectID(),
			URL:       url,
			Domain:    "amazon.co.uk",
			ItemName:  productName,
			Img:       imgSrc,
			Prices:    []model.Price{{ID: primitive.NewObjectID(), Price: price, CreatedAt: time.Now()}},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err := repo.SearchesCollection.InsertOne(context.TODO(), search)
		if err != nil {
			return nil, err
		}
		return &search, nil
	} else if err != nil {
		return nil, err
	}

	lastPrice := search.Prices[len(search.Prices)-1].Price
	if lastPrice != price {
		newPrice := model.Price{ID: primitive.NewObjectID(), Price: price, CreatedAt: time.Now()}
		search.Prices = append(search.Prices, newPrice)
		search.UpdatedAt = time.Now()
		_, err := repo.SearchesCollection.UpdateOne(
			context.TODO(),
			bson.M{"_id": search.ID},
			bson.M{"$set": bson.M{"prices": search.Prices, "updatedAt": search.UpdatedAt}},
		)
		if err != nil {
			return nil, err
		}
	}

	return &search, nil

}
