package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/fullstack/dev-overflow/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const QUESTIONCOLL = "questions"

type Dropper interface {
	Drop(context.Context) error
}

type MongoQuestionStore struct {
	client *mongo.Client
	coll *mongo.Collection
	TagStore
	UserStore
}

func NewMongoQuestionStore(client *mongo.Client, tagStore TagStore, userStore UserStore) *MongoQuestionStore {
	var mongoenvdbname = os.Getenv("MONGO_DB_NAME")
	return &MongoQuestionStore{
		client: client,
		coll: client.Database(mongoenvdbname).Collection(QUESTIONCOLL),
		TagStore: tagStore,
		UserStore: userStore,
	}
}

type QuestionStore interface {
	Dropper
	GetQuestionByID(context.Context, string) (*types.Question, error)
	GetQuestionsByUserID(context.Context, string) ([]*types.Question, error)
	GetQuestions(context.Context) ([]*types.Question, error)
	// TODO: IMPLEMENT SAVED QUESTIONS PAKAI PARAMS QUERY *types.SavedQuestionQueryParams
	GetQuestionsByTagID(context.Context, string) ([]*types.Question, error)
	GetSavedQuestions(context.Context, string) ([]*types.Question, error)
	AskQuestion(context.Context, *types.Question) (*types.Question, error)
	UpvoteQuestion(context.Context, *types.QuestionVoteParams) error
	DownvoteQuestion(context.Context, *types.QuestionVoteParams) error
	UpdateQuestionViews(context.Context, string) error
	UpdateQuestionAnswersField(context.Context, *types.UpdateQuestionAnswersParams) error
	DeleteQuestionByID(context.Context, string) error
	DeleteManyQuestionsByUserID(context.Context, primitive.ObjectID) error
}

func (s *MongoQuestionStore) Drop(ctx context.Context) error {
	fmt.Println("****DELETING DATABASE****")
	return s.coll.Drop(ctx)
}

func (s *MongoQuestionStore) GetQuestionByID(ctx context.Context, id string) (*types.Question, error) {
	var question types.Question

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": oid},
		},
		{
			"$lookup": bson.M{
				"from": "users",
				"localField": "userID",
				"foreignField": "_id",
				"as": "user",
			},
		},
		{
			"$unwind": "$user",
		},
		{
			"$lookup": bson.M{
				"from": "tags",
				"localField": "tags",
				"foreignField": "_id",
				"as": "tagDetails",	
			},
		},
	}

	cursor, err := s.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if cursor.Next(ctx) {
		if err := cursor.Decode(&question); err != nil {
			return nil, err
		}
	} else {
		return nil, mongo.ErrNoDocuments
	}

	return &question, nil
}

func (s *MongoQuestionStore) GetQuestionsByUserID(ctx context.Context, id string) ([]*types.Question, error) {
	var questions []*types.Question

	user, err := s.UserStore.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	oid := user.ID

	pipeline := []bson.M{
		{
			"$match": bson.M{"_id": oid},
		},
		{
			"$lookup": bson.M{
				"from": "users",
				"localField": "userID",
				"foreignField": "_id",
				"as": "user",
			},
		},
		{
			"$unwind": "$user",
		},
		{
			"$lookup": bson.M{
				"from": "tags",
				"localField": "tags",
				"foreignField": "_id",
				"as": "tagDetails",	
			},
		},
		{
			"$sort": bson.M{"upvotes": -1, "views": -1},
		},
	}

	cursor, err := s.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var question types.Question
		if err := cursor.Decode(&question); err != nil {
			return nil, err
		}

		questions = append(questions, &question)
	}

	return questions, nil

}

func (s *MongoQuestionStore) GetQuestions(ctx context.Context) ([]*types.Question, error) {
	var questions []*types.Question

	pipeline := []bson.M{
		{
			"$lookup": bson.M{
				"from": "users",
				"localField": "userID",
				"foreignField": "_id",
				"as": "user",
			}},
			{"$unwind":"$user"},
			{"$lookup":bson.M{
				"from": "tags",
				"localField": "tags",
				"foreignField": "_id",
				"as": "tagDetails",
			}},
			{"$sort":bson.M{"createdAt":-1}},
	}

	cursor, err := s.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var question types.Question
		if err := cursor.Decode(&question); err != nil {
			return nil, err
		}

		questions = append(questions, &question)
	}

	log.Println(questions)
	return questions, nil
}

func (s *MongoQuestionStore) GetSavedQuestions(ctx context.Context, id string) ([]*types.Question, error) {
	var questions []*types.Question
	
	user, err := s.UserStore.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, question := range user.Saved {
		pipeline := []bson.M{
			{
				"$match": bson.M{"_id": question},
			},
			// {
			// 	"$match": bson.M{"title": bson.M{"$regex": primitive.Regex{Pattern: params.SearchQuery, Options: "i"}}},
			// },
			{
				"$lookup": bson.M{
					"from": "users",
					"localField": "userID",
					"foreignField": "_id",
					"as": "user",
				},
			},
			{
				"$unwind": "$user",
			},
			{
				"$lookup": bson.M{
					"from": "tags",
					"localField": "tags",
					"foreignField": "_id",
					"as": "tagDetails",
				},
			},
			{"$sort":bson.M{"createdAt":-1}},
		}

		cursor, err := s.coll.Aggregate(ctx, pipeline)
		if err != nil {
			return nil, err
		}

		var question types.Question
		if cursor.Next(ctx) {
			if err := cursor.Decode(&question); err != nil {
				return nil, err
			}
			questions = append(questions, &question)
		} 
	}

	// startIndex := (params.Page - 1) * params.Limit
	// endIndex := params.Page + params.Limit

	// if startIndex > int64(len(questions)) {
	// 	startIndex = int64(len(questions))
	// }

	// if endIndex > int64(len(questions)) {
	// 	endIndex = int64(len(questions))
	// }

	// pagedQuestions := questions[startIndex:endIndex]

	// return pagedQuestions, nil
	return questions, nil
}

func (s *MongoQuestionStore) GetQuestionsByTagID(ctx context.Context, id string) ([]*types.Question, error) {
	var questions []*types.Question

	tag, err := s.TagStore.GetTagByID(ctx, id)
	if err != nil {
		return nil, err
	}

	pipeline := []bson.M{
		{"$match": bson.M{"tags": tag.ID}},
		{
			"$lookup": bson.M{
				"from": "users",
				"localField": "userID",
				"foreignField": "_id",
				"as": "user",
			}},
		{"$unwind": "$user"},
		{"$lookup": bson.M{
			"from": "tags",
			"localField": "tags",
			"foreignField": "_id",
			"as": "tagDetails",
		}},
		{"$sort": bson.M{"createdAt": -1}},
	}

	cursor, err := s.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var question types.Question
		if err := cursor.Decode(&question); err != nil {
			return nil, err
		}

		questions = append(questions, &question)
	}

	return questions, nil
}

func (s *MongoQuestionStore) AskQuestion(ctx context.Context, question *types.Question) (*types.Question, error) {
	res, err := s.coll.InsertOne(ctx, question)
	if err != nil {
		return nil, err
	}

	question.ID = res.InsertedID.(primitive.ObjectID)

	for _, tag := range question.Tags {
		if err := s.TagStore.UpdateTag(ctx, Map{"_id": tag}, &types.UpdateTagQuestionAndFollowers{Questions: question.ID, Followers: question.UserID}); err != nil {
			return nil, err
		}
	}

	_ = s.UserStore.UpdateUserQuestionsField(ctx, question.UserID, question.ID)

	return question, nil
}

func (s *MongoQuestionStore) UpvoteQuestion(ctx context.Context, params *types.QuestionVoteParams) error {
	user, err := s.UserStore.GetUserByID(ctx, params.UserID)
	if err != nil {
		return err
	}

	question, err := s.GetQuestionByID(ctx, params.QuestionID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": question.ID}

	updateDoc := bson.M{
		"$pull": bson.M{"downvotes": user.ID},
		"$addToSet": bson.M{"upvotes": user.ID},
	}

	_, err = s.coll.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoQuestionStore) DownvoteQuestion(ctx context.Context, params *types.QuestionVoteParams) error {
	user, err := s.UserStore.GetUserByID(ctx, params.UserID)
	if err != nil {
		return err
	}

	question, err := s.GetQuestionByID(ctx, params.QuestionID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": question.ID}

	updateDoc := bson.M{
		"$pull": bson.M{"upvotes": user.ID},
		"$addToSet": bson.M{"downvotes": user.ID},
	}

	_, err = s.coll.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}

	return nil

}

func (s *MongoQuestionStore) UpdateQuestionViews(ctx context.Context, id string) error {
	question, err := s.GetQuestionByID(ctx, id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": question.ID}

	_ = s.coll.FindOneAndUpdate(ctx, filter, bson.M{"$inc": bson.M{"views": 1}})

	return nil

}

func (s *MongoQuestionStore) UpdateQuestionAnswersField(ctx context.Context, update *types.UpdateQuestionAnswersParams) error {
	question, err := s.GetQuestionByID(ctx, update.QuestionID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": question.ID}

	updateDoc := bson.M{
		"$push": bson.M{
			"answers": bson.M{"$each":[]primitive.ObjectID{update.Answers}},
		},
	}

	_, err = s.coll.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoQuestionStore) DeleteQuestionByID(ctx context.Context, id string) error {

	oid,err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	if err := s.coll.FindOneAndDelete(ctx, bson.M{"_id": oid}).Err(); err != nil {
		return err
	}
	return nil
}

func (s *MongoQuestionStore) DeleteManyQuestionsByUserID(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.coll.DeleteMany(ctx, bson.M{"userID": id})
	if err != nil {
		return err
	}
	return nil
}