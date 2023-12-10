package api

import (
	"time"

	"github.com/fullstack/dev-overflow/db"
	"github.com/fullstack/dev-overflow/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AnswerHandler struct {
	answerStore db.AnswerStore
	questionStore db.QuestionStore
}

func NewAnswerHandler(answerStore db.AnswerStore, questionStore db.QuestionStore) *AnswerHandler {
	return &AnswerHandler{
		answerStore: answerStore,
		questionStore: questionStore,
	}
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
		Content: params.Content,
		Upvotes: []primitive.ObjectID{},
		Downvotes: []primitive.ObjectID{},
		CreatedAt: time.Now().UTC(),
	}

	answer, err := h.answerStore.CreateAnswer(ctx.Context(), answer)
	if err != nil {
		return ErrBadRequest()
	}

	// if err := h.questionStore.UpdateQuestionAnswersField(ctx.Context(), db.Map{"_id":params.QuestionID}, &types.UpdateQuestionAnswersParams{Answers: answer.ID}); err != nil {
	// 	return ErrBadRequest()
	// }

	return ctx.JSON(answer)
}