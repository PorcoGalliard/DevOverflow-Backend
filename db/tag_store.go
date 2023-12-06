package db

import (
	"context"
	"errors"
	"os"

	"github.com/fullstack/dev-overflow/types"
	"github.com/fullstack/dev-overflow/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const TAGCOLL = "tags"

type MongoTagStore struct {
	client *mongo.Client
	collection *mongo.Collection
}

func NewMongoTagStore(client *mongo.Client) *MongoTagStore {
	var mongoEnvDBName = os.Getenv("MONGO_DB_NAME")
	return &MongoTagStore{
		client: client,
		collection: client.Database(mongoEnvDBName).Collection(TAGCOLL),
	}
}

type TagStore interface {
	CreateTag(context.Context, *types.Tag) (*types.Tag, error)
	GetTagByID(context.Context, string) (*types.Tag, error)
	GetTagByName(context.Context, string) (*types.Tag, error)
	GetTags(context.Context) ([]*types.Tag, error)
	UpdateTag(context.Context, Map, *types.UpdateTagQuestionAndFollowers) error
	UpdateManyFollowersByID(context.Context, primitive.ObjectID) error
	UpdateManyQuestionsByID(context.Context, primitive.ObjectID) error
}

func (s *MongoTagStore) GetTagByID(ctx context.Context, id string) (*types.Tag, error) {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var tag types.Tag
	if err := s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&tag); err != nil {
		return nil, err
	}

	return &tag, nil
}

func (s *MongoTagStore) GetTagByName(ctx context.Context, name string) (*types.Tag, error) {
	var tag types.Tag
	name = utils.FormatTag(name)
	if err := s.collection.FindOne(ctx, bson.M{"name": name}).Decode(&tag); err != nil {
		return nil, err
	}

	return &tag, nil
}

func (s *MongoTagStore) GetTags(ctx context.Context) ([]*types.Tag, error) {
	var tags []*types.Tag
	cursor, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &tags); err != nil {
		return nil, err
	}

	return tags, nil

}

func (s *MongoTagStore) CreateTag(c context.Context, tag *types.Tag) (*types.Tag, error) {
	tag.Name = utils.FormatTag(tag.Name)
	res, err := s.collection.InsertOne(c, tag)
	if err != nil {
		return nil , err
	}
	
	tag.ID = res.InsertedID.(primitive.ObjectID)

	return tag, nil
}

func (s *MongoTagStore) UpdateTag(ctx context.Context, filter Map, update *types.UpdateTagQuestionAndFollowers) error {

	oid, ok := filter["_id"]
	if !ok {
		return errors.New("filter[_id] is not a primitive.ObjectID")
	}

	filter["_id"] = oid

	updateDoc := bson.M{
		"$push": bson.M{
			"questions": bson.M{"$each":[]primitive.ObjectID{update.Questions}},
			"followers": bson.M{"$each":[]primitive.ObjectID{update.Followers}},
		},
	}

	_, err := s.collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoTagStore) UpdateManyFollowersByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.collection.UpdateMany(ctx, bson.M{"followers": id}, bson.M{"$pull": bson.M{"followers": id}})
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoTagStore) UpdateManyQuestionsByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.collection.UpdateMany(ctx, bson.M{"questions": id}, bson.M{"$pull": bson.M{"questions": id}})
	if err != nil {
		return err
	}

	return nil
}