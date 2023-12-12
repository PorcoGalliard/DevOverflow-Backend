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
	GetAnswerByID(context.Context, string) (*types.Answer, error)
	GetAnswersOfQuestion(context.Context, string) ([]*types.Answer, error)
	CreateAnswer(context.Context, *types.Answer) (*types.Answer,error)
	UpvoteAnswer(context.Context, *types.VoteAnswerParams) error
	DownvoteAnswer(context.Context, *types.VoteAnswerParams) error
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

func (s *MongoAnswerStore) GetAnswerByID(ctx context.Context, id string) (*types.Answer, error) {
	var answer types.Answer

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	if err := s.coll.FindOne(ctx, Map{"_id": oid}).Decode(&answer); err != nil {
		return nil, err
	}

	return &answer, nil
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

func (s *MongoAnswerStore) UpvoteAnswer(ctx context.Context, params *types.VoteAnswerParams) error {
	answer, err := s.GetAnswerByID(ctx, params.AnswerID)
	if err != nil {
		return err
	}

	userID, err := primitive.ObjectIDFromHex(params.UserID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": answer.ID}

	updateDoc := bson.M{
		"$pull":bson.M{"downvotes": userID},
		"$addToSet":bson.M{"upvotes": userID},
	}

	_, err = s.coll.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoAnswerStore) DownvoteAnswer(ctx context.Context, params *types.VoteAnswerParams) error {
	answer, err := s.GetAnswerByID(ctx, params.AnswerID)
	if err != nil {
		return err
	}

	userID, err := primitive.ObjectIDFromHex(params.UserID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": answer.ID}

	updateDoc := bson.M{
		"$pull":bson.M{"upvotes": userID},
		"$addToSet":bson.M{"downvotes": userID},
	}

	_, err = s.coll.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}

	return nil
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