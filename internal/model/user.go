package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Username       string               `bson:"username" json:"username"`
	Email          string               `bson:"email" json:"email"`
	HashedPassword string               `bson:"hashedPassword" json:"hashedPassword"`
	MySearches     []primitive.ObjectID `bson:"mySearches,omitempty" json:"mySearches,omitempty"`
	CreatedAt      time.Time            `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt      time.Time            `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}
