package api

import (
	"errors"
	"log"
	"time"

	"github.com/fullstack/dev-overflow/db"
	"github.com/fullstack/dev-overflow/types"
	"github.com/fullstack/dev-overflow/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuestionHandler struct {
	questionStore db.QuestionStore
	userStore db.UserStore
	tagStore db.TagStore
	answerStore db.AnswerStore
}

func NewQuestionHandler(questionStore db.QuestionStore, userStore db.UserStore, tagStore db.TagStore, answerStore db.AnswerStore) *QuestionHandler {
	return &QuestionHandler{
		questionStore: questionStore,
		userStore: userStore,
		tagStore: tagStore,
		answerStore: answerStore,
	}
}

func (h *QuestionHandler) HandleGetQuestionByID(ctx *fiber.Ctx) error {
	var (
		id = ctx.Params("id")
	)
	
	question, err := h.questionStore.GetQuestionByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound(id)
		}
		return err
	}

	return ctx.JSON(question)
}

func (h *QuestionHandler) HandleGetQuestionsByUserID(ctx *fiber.Ctx) error {
	var (
		id = ctx.Params("id")
	)

	questions, err := h.questionStore.GetQuestionsByUserID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound(id)
		}
		return err
	}

	return ctx.JSON(questions)
}

func (h *QuestionHandler) HandleGetQuestions(ctx *fiber.Ctx) error {
	questions, err := h.questionStore.GetQuestions(ctx.Context())
	if err != nil {
		return ErrResourceNotFound("question")
	}

	return ctx.JSON(questions)
}

func (h *QuestionHandler) HandleGetSavedQuestions(ctx *fiber.Ctx) error {
	var (
		id = ctx.Params("clerkID")
		// params types.SavedQuestionQueryParams
	)

	// if err := ctx.BodyParser(&params); err != nil {
	// 	return ErrBadRequest()
	// }

	// questions, err := h.questionStore.GetSavedQuestions(ctx.Context(), id, &params)
	questions, err := h.questionStore.GetSavedQuestions(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound(id)
		}
		return err
	}

	return ctx.JSON(questions)
}

func (h *QuestionHandler) HandleGetQuestiosByTagID (ctx *fiber.Ctx) error {
	var (
		id = ctx.Params("id")
	)

	questions, err := h.questionStore.GetQuestionsByTagID(ctx.Context(), id)
	if err != nil {
		return ErrResourceNotFound(id)
	}

	return ctx.JSON(questions)
}

func (h *QuestionHandler) HandleAskQuestion(ctx *fiber.Ctx) error {
	var params types.AskQuestionParams

	if err := ctx.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	if errors := params.Validate(); len(errors) > 0 {
		return ctx.JSON(errors)
	}

	for i, tag := range params.Tags {
		tag = utils.FormatTag(tag)
		params.Tags[i] = tag
	}

	user, err := h.userStore.GetUserByID(ctx.Context(), params.ClerkID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound(params.ClerkID)
		}
	}

	tags := make([]primitive.ObjectID, len(params.Tags))
	for i, tagName := range params.Tags {

		tag, err := h.tagStore.GetTagByName(ctx.Context(), tagName)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				tag = &types.Tag{
					Name: tagName,
					Questions: []primitive.ObjectID{},
					Followers: []primitive.ObjectID{},
					CreatedAt: time.Now().UTC(),
				}

				insertedTag, err := h.tagStore.CreateTag(ctx.Context(), tag)
				if err != nil {
					return ErrBadRequest()
				}

				tag = insertedTag
			} 
		}
		tags[i] = *&tag.ID
	}

	question := &types.Question{
		Title: params.Title,
		Description: params.Description,
		UserID: user.ID,
		Tags: tags,
		Upvotes: []primitive.ObjectID{},
		Downvotes: []primitive.ObjectID{},
		Answers: []primitive.ObjectID{},
		CreatedAt: time.Now().UTC(),
	}


	insertedQuestion, err := h.questionStore.AskQuestion(ctx.Context(), question)
	if err != nil {
		log.Println(err)
			return ErrBadRequest()
		}

	return ctx.JSON(insertedQuestion)
}

func (h *QuestionHandler) HandleDeleteQuestionByID(ctx *fiber.Ctx) error {
	var (
		id = ctx.Params("_id")
		params types.DeleteQuestionParams
	)

	if err := ctx.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	user, err := h.userStore.GetUserByID(ctx.Context(), params.UserID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound(params.UserID)
		}
	}

	question, err := h.questionStore.GetQuestionByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound(id)
		}
	}

	for _, answerID := range question.Answers {
		if err := h.answerStore.DeleteAnswerByID(ctx.Context(), answerID.Hex()); err != nil {
			return ErrBadRequest()
		}
	}

	if question.UserID != user.ID {
		return ErrUnauthorized()
	}

	if question.UserID == user.ID {
		if err := h.questionStore.DeleteQuestionByID(ctx.Context(), id); err != nil {
			return ErrBadRequest()
		
		}
	}

	return nil
}

func (h *QuestionHandler) HandleQuestionVote(ctx *fiber.Ctx) error {
	var params types.QuestionVoteParams

	if err := ctx.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	question, err := h.questionStore.GetQuestionByID(ctx.Context(), params.QuestionID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound(params.QuestionID)
		}
	}

	user, err := h.userStore.GetUserByID(ctx.Context(), params.UserID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound(params.UserID)
		}
	}

	if question.UserID == user.ID {
		return ErrUnauthorized()
	}

	if !params.HasUpvoted && !params.HasDownvoted {
		return ErrBadRequest()
	}

	if params.HasUpvoted {
		if err := h.questionStore.UpvoteQuestion(ctx.Context(), &params); err != nil {
			return ErrBadRequest()
		}
	}

	if params.HasDownvoted {
		if err := h.questionStore.DownvoteQuestion(ctx.Context(), &params); err != nil {
			return ErrBadRequest()
		}
	}

	question, err = h.questionStore.GetQuestionByID(ctx.Context(), params.QuestionID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound(params.QuestionID)
		}
	}

	return ctx.JSON(question)
}