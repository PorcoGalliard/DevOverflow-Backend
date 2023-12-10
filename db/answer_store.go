package db

import (
	"context"
	"os"

	"github.com/fullstack/dev-overflow/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const ANSWERCOLL = "answers"

type AnswerStore interface {
	GetAnswersOfQuestion(context.Context, string) ([]*types.Answer, error)
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

func (s *MongoAnswerStore) GetAnswersOfQuestion(ctx context.Context, id string) ([]*types.Answer, error) {
	var answers []*types.Answer

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{"questionID": oid},
		},
		{
			"$lookup":bson.M{
				"from": "users",
				"localField": "userID",
				"foreignField": "_id",
				"as": "user",
			},
		},
		{
			"$unwind": "$user",
		},
		{"$sort": bson.M{"createdAt": -1}},
	}

	cursor, err := s.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &answers); err != nil {
		return nil, err
	}

	return answers, nil
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