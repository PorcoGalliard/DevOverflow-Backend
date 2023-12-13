package db

import (
	"context"
	"os"

	"github.com/fullstack/dev-overflow/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const USERCOLL = "users"

type MongoUserStore struct {
	client *mongo.Client
	coll *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	var mongoEnvDBName = os.Getenv("MONGO_DB_NAME")
	return &MongoUserStore{
		client: client,
		coll: client.Database(mongoEnvDBName).Collection(USERCOLL),
	}
}

type UserStore interface {
	CreateUser(context.Context, *types.User) (*types.User, error)
	GetUserByID(context.Context, string) (*types.User, error)
	GetUsers(context.Context, UserQueryParams) ([]*types.User, error)
	SaveQuestion(context.Context, *types.SaveQuestionParam) (bool  ,error)
	UpdateUserQuestionsField(context.Context, primitive.ObjectID, primitive.ObjectID) error
	UpdateUserAnswersField(context.Context, primitive.ObjectID, primitive.ObjectID) error
	UpdateUser(context.Context, string, *types.UpdateUserParam) error
	DeleteUser(context.Context, string) error
}

func (s *MongoUserStore) CreateUser(c context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(c, user)
	if err != nil {
		return nil , err
	}
	
	user.ID = res.InsertedID.(primitive.ObjectID)

	return user, nil
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"clerkID": id}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context, params UserQueryParams) ([]*types.User, error) {
	var users []*types.User

	opt := options.FindOptions{}
	opt.SetSkip((params.Page - 1) * params.Limit)
	opt.SetLimit(params.Limit)


	query := bson.M{}

	if params.SearchQuery != "" {
		query["$or"] = []bson.M{
			{"firstName": bson.M{"$regex":params.SearchQuery, "$options": "i"}},
			{"lastName": bson.M{"$regex":params.SearchQuery, "$options": "i"}},
		}
	}

	switch params.Filter {
	case "old_user":
		opt.SetSort(bson.M{"joinedAt": 1})
	case "new_user":
		opt.SetSort(bson.M{"joinedAt": -1})
	case "top_contributor":
		opt.SetSort(bson.M{"reputation": -1})
	}

	cursor, err := s.coll.Find(ctx, query, &opt)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, clerkID string, update *types.UpdateUserParam) error {

	filter := bson.M{"clerkID": clerkID}
	updateData := bson.M{"$set": update.UpdateData}

	result := s.coll.FindOneAndUpdate(ctx, filter, updateData)

	var updatedUser types.User
	if err := result.Decode(&updatedUser); err != nil {
		return err
	}

	return nil
}

func (s *MongoUserStore) UpdateUserQuestionsField(ctx context.Context, userID primitive.ObjectID, questionID primitive.ObjectID) error {
	filter := bson.M{"_id": userID}
	updateData := bson.M{
		"$push": bson.M{
			"questions": bson.M{"$each":[]primitive.ObjectID{questionID}},
		},
	}

	result := s.coll.FindOneAndUpdate(ctx, filter, updateData)

	var updatedUser types.User
	if err := result.Decode(&updatedUser); err != nil {
		return err
	}

	return nil
}

func (s *MongoUserStore) UpdateUserAnswersField(ctx context.Context, userID primitive.ObjectID, answerID primitive.ObjectID) error {
	filter := bson.M{"_id": userID}
	updateData := bson.M{
		"$push": bson.M{
			"answers": bson.M{"$each":[]primitive.ObjectID{answerID}},
		},
	}

	result := s.coll.FindOneAndUpdate(ctx, filter, updateData)

	var updatedUser types.User
	if err := result.Decode(&updatedUser); err != nil {
		return err
	}

	return nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, clerkID string) error {
	user, err := s.GetUserByID(ctx, clerkID)
	if err != nil {
		return err
	}

	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": user.ID})
	if err != nil {
		return err
	}
	
	return nil
}

func (s *MongoUserStore) SaveQuestion(ctx context.Context, params *types.SaveQuestionParam) (bool, error) {

	var updatedUser types.User

	user, err := s.GetUserByID(ctx, params.UserID)
	if err != nil {
		return false, err
	}

	paramsOID, err := primitive.ObjectIDFromHex(params.QuestionID)
	if err != nil {
		return false, err
	}

	for _, questionID := range user.Saved {
		if questionID == paramsOID {
			updateData := bson.M{"$pull": bson.M{"saved": paramsOID}}
			result := s.coll.FindOneAndUpdate(ctx, bson.M{"clerkID": params.UserID}, updateData)
			if err := result.Decode(&updatedUser); err != nil {
				return false, err
			}
			return false, nil
		}
	}

	filter := bson.M{"clerkID": params.UserID}
	updateData := bson.M{"$push": bson.M{"saved": paramsOID}}
	result := s.coll.FindOneAndUpdate(ctx, filter, updateData)
	if err := result.Decode(&updatedUser); err != nil {
		return false, err
	}

	return true, nil
}