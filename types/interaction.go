package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Interaction struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID primitive.ObjectID `bson:"userID" json:"userID"`
	Action string `bson:"action" json:"action"`
	QuestionID primitive.ObjectID `bson:"questionID" json:"questionID"`
	AnswerID primitive.ObjectID `bson:"answerID" json:"answerID"`
	Tags []primitive.ObjectID `bson:"tags" json:"tags"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type ViewQuestionParams struct {
	UserID string `json:"userID,omitempty"`
	QuestionID string `json:"questionID"`
}