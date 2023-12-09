package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Answer struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
	QuestionID primitive.ObjectID `bson:"question_id" json:"question_id"`
	Content string `bson:"content" json:"content"`
	Upvotes []primitive.ObjectID `bson:"upvotes,omitempty" json:"upvotes,omitempty"`
	Downvotes []primitive.ObjectID `bson:"downvotes" json:"downvotes"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type CreateAnswerParams struct {
	UserID primitive.ObjectID `json:"user_id"`
	QuestionID primitive.ObjectID `json:"question_id"`
	Content string `json:"content"`
}