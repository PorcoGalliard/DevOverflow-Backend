package fixtures

import (
	"context"
	"log"
	"time"

	"github.com/fullstack/dev-overflow/db"
	"github.com/fullstack/dev-overflow/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddQuestion(store *db.Store, title string, desc string, userID primitive.ObjectID, tags []primitive.ObjectID, createdAt time.Time) (*types.Question) {

	question := &types.Question{
		Title: title,
		Description: desc,
		UserID: userID,
		Tags: tags,
		Upvotes: []primitive.ObjectID{},
		Downvotes: []primitive.ObjectID{},
		Answers: []primitive.ObjectID{},
		CreatedAt: createdAt,
	}

	insertedQuestion, err := store.Question.AskQuestion(context.Background(), question)
	if err != nil {
		log.Fatal(err)
	}

	return insertedQuestion
}

func AddUser(store *db.Store, firstName string, lastName string, clerkID string, email string, encryptedPassword string) (*types.User) {
	user, err := types.NewUserFromParams(types.CreateUserParam{
		FirstName: firstName,
		LastName: lastName,
		ClerkID: clerkID,
		Email: email,
		// Password: encryptedPassword,
	})

	if err != nil {
		log.Fatal(err)
	}

	insertedUser, err := store.User.CreateUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}

	return insertedUser
}

func AddTag(store *db.Store, name string) (*types.Tag) {
	tag := &types.Tag{
		Name: name,
		Questions: []primitive.ObjectID{},
		Followers: []primitive.ObjectID{},
		CreatedAt: time.Now().UTC(),
	}

	insertedTag, err := store.Tag.CreateTag(context.Background(), tag)
	if err != nil {
		log.Fatal(err)
	}

	return insertedTag
}

func UpdateTag(store *db.Store, tagID primitive.ObjectID, update *types.UpdateTagQuestionAndFollowers) (*types.Tag, error) {
	if err := store.Tag.UpdateTag(context.Background(), db.Map{"_id":tagID}, update); err != nil {
		log.Fatal(err)
	}
	
	realTag, err := store.Tag.GetTagByID(context.Background(), tagID.Hex())
	if err != nil {
		log.Fatal(err)
	}

	return realTag, nil
}