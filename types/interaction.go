package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Interaction struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID primitive.ObjectID `bson:"userID" json:"userID"`
	User *User `bson:"user,omitempty" json:"user,omitempty"`
	Action string `bson:"action" json:"action"`
	QuestionID primitive.ObjectID `bson:"questionID" json:"questionID"`
	AnswerID primitive.ObjectID `bson:"answerID" json:"answerID"`
	Tags []primitive.ObjectID `bson:"tags" json:"tags"`
	TagsDetails []*Tag `bson:"tagsDetails,omitempty" json:"tagsDetails,omitempty"`
	CreatedAt int64 `bson:"createdAt" json:"createdAt"`
}

type ViewQuestionParams struct {
	UserID string `json:"userID,omitempty"`
	QuestionID string `json:"questionID"`
}