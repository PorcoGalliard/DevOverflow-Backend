package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tag struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Description string `bson:"description" json:"description"`
	Name string `bson:"name" json:"name"`
	Questions []primitive.ObjectID `bson:"questions" json:"questions"`
	QuestionDetails []*Question `bson:"questionDetails,omitempty" json:"questionDetails,omitempty"`
	Followers []primitive.ObjectID `bson:"followers" json:"followers"`
	FollowersDetails []*User `bson:"followersDetails,omitempty" json:"followersDetails,omitempty"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

type CreateTagParams struct {
	Name string `json:"name"`
}

type UpdateTagQuestionAndFollowers struct {
	Questions primitive.ObjectID `json:"questions"`
	Followers primitive.ObjectID `json:"followers"`
}