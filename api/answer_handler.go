package api

import (
	"errors"
	"time"

	"github.com/fullstack/dev-overflow/db"
	"github.com/fullstack/dev-overflow/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AnswerHandler struct {
	answerStore db.AnswerStore
	questionStore db.QuestionStore
	userStore db.UserStore
}

func NewAnswerHandler(answerStore db.AnswerStore, questionStore db.QuestionStore, userStore db.UserStore) *AnswerHandler {
	return &AnswerHandler{
		answerStore: answerStore,
		questionStore: questionStore,
		userStore: userStore,
	}
}

func (h *AnswerHandler) HandleGetAnswerByID(ctx *fiber.Ctx) error {
	var (
		questionID = ctx.Params("questionID")
		answerID = ctx.Params("answerID")
	)

	question, err := h.questionStore.GetQuestionByID(ctx.Context(), questionID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound(questionID)
		}
	}

	answer, err := h.answerStore.GetAnswerByID(ctx.Context(), answerID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound(answerID)
		}
	}

	if answer.QuestionID != question.ID {
		return ErrBadRequest()
	}

	return ctx.JSON(answer)
}

func (h *AnswerHandler) HandleAnswerVote(ctx *fiber.Ctx) error {
	var params types.VoteAnswerParams

	if err := ctx.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	answer, err := h.answerStore.GetAnswerByID(ctx.Context(), params.AnswerID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound(params.AnswerID)
		}
	}

	user, err := h.userStore.GetUserByID(ctx.Context(), params.UserID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound(params.UserID)
		}
	}

	if answer.UserID == user.ID {
		return ErrUnauthorized()
	}

	if !params.HasUpvoted && !params.HasDownvoted {
		return ErrBadRequest()
	}

	if params.HasUpvoted {
		if err := h.answerStore.UpvoteAnswer(ctx.Context(), &params); err != nil {
			return ErrBadRequest()
		}
	}

	if params.HasDownvoted {
		if err := h.answerStore.DownvoteAnswer(ctx.Context(), &params); err != nil {
			return ErrBadRequest()
		}
	}

	answer, err = h.answerStore.GetAnswerByID(ctx.Context(), params.AnswerID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound(params.AnswerID)
		}
	}
	return ctx.JSON(answer)
}


func (h *AnswerHandler) HandleGetAnswersOfQuestion(ctx *fiber.Ctx) error {

	var (
		id = ctx.Params("id")
	)

	answers, err := h.answerStore.GetAnswersOfQuestion(ctx.Context(), id)
	if err != nil {
		return ErrBadRequest()
	}

	return ctx.JSON(answers)
}

func (h *AnswerHandler) HandleCreateAnswer(ctx *fiber.Ctx) error {
	var (
		params types.CreateAnswerParams
	)

	if err := ctx.BodyParser(&params); err != nil {
		return err
	}

	answer := &types.Answer{
		UserID: params.UserID,
		QuestionID: params.QuestionID,
		Description: params.Description,
		Upvotes: []primitive.ObjectID{},
		Downvotes: []primitive.ObjectID{},
		CreatedAt: time.Now().UTC(),
	}

	answer, err := h.answerStore.CreateAnswer(ctx.Context(), answer)
	if err != nil {
		return ErrBadRequest()
	}

	if err := h.questionStore.UpdateQuestionAnswersField(ctx.Context(), &types.UpdateQuestionAnswersParams{
		QuestionID: params.QuestionID.Hex(),
		Answers: answer.ID,
	}); err != nil {
		return ErrBadRequest()
	}

	return ctx.JSON(answer)
}