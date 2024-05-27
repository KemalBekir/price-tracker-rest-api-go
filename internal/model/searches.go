package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Searches struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	URL       string             `bson:"url" json:"url"`
	Domain    string             `bson:"domain" json:"domain"`
	ItemName  string             `bson:"itemName" json:"itemName"`
	Img       string             `bson:"img,omitempty" json:"img,omitempty"`
	Prices    []Price            `bson:"prices,omitempty" json:"prices,omitempty"`
	Owner     primitive.ObjectID `bson:"owner,omitempty" json:"owner,omitempty"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
