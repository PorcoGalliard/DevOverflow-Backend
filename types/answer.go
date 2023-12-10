package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Answer struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID primitive.ObjectID `bson:"userID" json:"userID"`
	QuestionID primitive.ObjectID `bson:"questionID" json:"questionID"`
	Description string `bson:"content" json:"description"`
	Upvotes []primitive.ObjectID `bson:"upvotes,omitempty" json:"upvotes,omitempty"`
	Downvotes []primitive.ObjectID `bson:"downvotes" json:"downvotes"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type CreateAnswerParams struct {
	UserID primitive.ObjectID `json:"userID"`
	QuestionID primitive.ObjectID `json:"questionID"`
	Description string `json:"description"`
}

type DeleteAnswerParams struct {
	QuestionID primitive.ObjectID `json:"questionID"`
}