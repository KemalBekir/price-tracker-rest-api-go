package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Price struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Price     float64            `bson:"price" json:"price"`
	SearchID  primitive.ObjectID `bson:"search" json:"search"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
