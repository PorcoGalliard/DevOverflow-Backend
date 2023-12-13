package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TagWithSingleQuestion struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Description string `bson:"description" json:"description"`
	Name string `bson:"name" json:"name"`
	Questions []primitive.ObjectID `bson:"questions" json:"questions"`
	Followers []primitive.ObjectID `bson:"followers" json:"followers"`
	FollowersDetails []*User `bson:"followersDetails,omitempty" json:"followersDetails,omitempty"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	QuestionDetails *Question `bson:"questionDetails,omitempty" json:"questionDetails,omitempty"`
}