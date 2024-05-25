package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Searches struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	URL         string             `bson:"url" json:"url"`
	Domain      string             `bson:"domain" json:"domain,omitempty"`
	ItemName    string             `bson:"itemName" json:"itemName"`
	Img         string             `bson:"img,omitempty" json:"img,omitempty"`
	Prices      []Price            `bson:"prices,omitempty" json:"prices,omitempty"`
	LatestPrice float64            `bson:"latestPrice,omitempty" json:"latestPrice,omitempty"`
	Owner       *User              `bson:"owner,omitempty" json:"owner,omitempty"`
	CreatedAt   time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt   time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
