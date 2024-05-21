package services

import (
	"context"

	"github.com/KemalBekir/price-tracker-rest-api-go/internal/db"
	"github.com/KemalBekir/price-tracker-rest-api-go/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
