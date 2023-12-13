package api

import (
	"errors"
	"time"

	"github.com/fullstack/dev-overflow/db"
	"github.com/fullstack/dev-overflow/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

type TagHandler struct {
	tagStore db.TagStore
	userStore db.UserStore
}

func NewTagHandler(tagStore db.TagStore, userStore db.UserStore) *TagHandler {
	return &TagHandler{
		tagStore: tagStore,
		userStore: userStore,
	}
}

func (h *TagHandler) HandleGetTagByID(ctx *fiber.Ctx) error {
	var (
		id = ctx.Params("_id")
	)

	tag, err := h.tagStore.GetTagByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments){
			return ErrResourceNotFound(id)
		}
	}
	return ctx.JSON(tag)

}

func (h *TagHandler) HandleGetTagByName(ctx *fiber.Ctx) error {
	var (
		name = ctx.Params("name")
	)

	tag, err := h.tagStore.GetTagByName(ctx.Context(), name)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments){
			return ErrResourceNotFound(name)
		}
	}
	return ctx.JSON(tag)
}

func (h *TagHandler) HandleGetTags(ctx *fiber.Ctx) error {
	tags, err := h.tagStore.GetTags(ctx.Context())
	if err != nil {
		return ErrBadRequest()
	}

	return ctx.JSON(tags)
}

func (h *TagHandler) HandleGetQuestionByTagID(ctx *fiber.Ctx) error {
	var (
		id = ctx.Params("id")
	)

	questions, err := h.tagStore.GetQuestionByTagID(ctx.Context(), id)
	if err != nil {
		return ErrBadRequest()
	}

	return ctx.JSON(questions)
}

func (h *TagHandler) HandleCreateTag(ctx *fiber.Ctx) error {
	var params types.CreateTagParams

	if err := ctx.BodyParser(&params); err != nil {
		return err
	}

	tag := &types.Tag{
		Name: params.Name,
		CreatedAt: time.Now().UTC(),
	}

	tag, err := h.tagStore.CreateTag(ctx.Context(), tag)
	if err != nil {
		return ErrBadRequest()
	}

	return ctx.JSON(tag)
}

func (h *TagHandler) HandleUpdateTag(ctx *fiber.Ctx) error {
	var (
		params *types.UpdateTagQuestionAndFollowers
		oid = ctx.Params("_id")
	)
	if err := ctx.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	filter := db.Map{"_id": oid}
	if err := h.tagStore.UpdateTag(ctx.Context(), filter, params); err != nil {
		return ErrBadRequest()
	}

	return ctx.JSON(fiber.Map{"message": "Tag updated successfully"})
}