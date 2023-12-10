package db

import (
	"context"
	"os"

	"github.com/fullstack/dev-overflow/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const ANSWERCOLL = "answers"

type AnswerStore interface {
	CreateAnswer(context.Context, *types.Answer) (*types.Answer,error)
	DeleteAnswerByID(context.Context, string) error
}

type MongoAnswerStore struct {
	client *mongo.Client
	coll *mongo.Collection
}

func NewMongoAnswerStore(client *mongo.Client) *MongoAnswerStore {
	var mongoenvdbname = os.Getenv("MONGO_DB_NAME")
	return &MongoAnswerStore{
		client: client,
		coll: client.Database(mongoenvdbname).Collection(ANSWERCOLL),
	}
}

func (s *MongoAnswerStore) CreateAnswer(ctx context.Context, answer *types.Answer) ( *types.Answer, error) {
	res, err := s.coll.InsertOne(ctx, answer)
	if err != nil {
		return nil, err
	}

	answer.ID = res.InsertedID.(primitive.ObjectID)

	return answer, nil
}

func (s *MongoAnswerStore) DeleteAnswerByID(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.coll.DeleteOne(ctx, Map{"_id": oid})
	if err != nil {
		return err
	}

	return nil
}