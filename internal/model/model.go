package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Searches struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	URL       string             `bson:"url" json:"url"`
	DOMAIN    string             `bson:"domain" json:"domain"`
	ITEM_NAME string             `bson:"itemName" json:"itemName"`
	IMG       string             `bson:"img" json:"img"`
	PRICES    primitive.ObjectID `bson:"prices" json:"prices"`
	OWNER     primitive.ObjectID `bson:"owner,omitempty" json:"owner,omitempty"`
}
