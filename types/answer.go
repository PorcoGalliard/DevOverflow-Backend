package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Answer struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
}