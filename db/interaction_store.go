package db

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/fullstack/dev-overflow/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const INTERACTIONCOLL = "interactions"

type InteractionStore interface {
	GetInteractionByUserAndQuestionID(context.Context, string, string) (*types.Interaction, error)
	CreateViewInteraction(context.Context, *types.ViewQuestionParams) (*types.Interaction, error)
}

type MongoInteractionStore struct {
	client *mongo.Client
	coll *mongo.Collection
	userStore UserStore
	questionStore QuestionStore
}

func NewMongoInteractionStore(client *mongo.Client, userStore UserStore, questionStore QuestionStore) *MongoInteractionStore {
	var mongoenvdbname = os.Getenv("MONGO_DB_NAME")
	return &MongoInteractionStore{
		client: client,
		coll: client.Database(mongoenvdbname).Collection(INTERACTIONCOLL),
		userStore: userStore,
		questionStore: questionStore,
	}
}

func (s *MongoInteractionStore) GetInteractionByUserAndQuestionID(ctx context.Context, userID string, questionID string) (*types.Interaction, error) {
	var interaction types.Interaction

	user, err := s.userStore.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	question, err := s.questionStore.GetQuestionByID(ctx, questionID)
	if err != nil {
		return nil, err
	}

	err = s.coll.FindOne(ctx, bson.M{"userID": user.ID, "questionID": question.ID}).Decode(&interaction); 
	if err != nil {
		if err == mongo.ErrNoDocuments{
			return nil, nil
		}
		return nil, err
	}

	return &interaction, nil
}

func (s *MongoInteractionStore) CreateViewInteraction(ctx context.Context, params *types.ViewQuestionParams) (*types.Interaction, error) {
	var interaction types.Interaction

	user, err := s.userStore.GetUserByID(ctx, params.UserID)
	if err != nil {
		return nil, err
	}

	question, err := s.questionStore.GetQuestionByID(ctx, params.QuestionID)
	if err != nil {
		return nil, err
	}

	existingInteraction, err := s.GetInteractionByUserAndQuestionID(ctx, params.UserID, params.QuestionID)
	if err != nil {
		return nil, err
	}

	if existingInteraction != nil {
		return nil, errors.New("user already viewed this question")
	}

	interaction = types.Interaction{
		UserID: user.ID,
		Action: "view",
		QuestionID: question.ID,
		CreatedAt: time.Now().UTC(),
	}

	res, err := s.coll.InsertOne(ctx, interaction)
	if err != nil {
		return nil, err
	}

	interaction.ID = res.InsertedID.(primitive.ObjectID)

	_ = s.questionStore.UpdateQuestionViews(ctx, params.QuestionID)

	return &interaction, nil
}