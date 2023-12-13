package api

import (
	"github.com/fullstack/dev-overflow/db"
	"github.com/fullstack/dev-overflow/types"
	"github.com/gofiber/fiber/v2"
)

type InteractionHandler struct {
	interactionStore db.InteractionStore
}

func NewInteractionHandler(interactionStore db.InteractionStore) *InteractionHandler {
	return &InteractionHandler{
		interactionStore: interactionStore,
	}
}

func (h *InteractionHandler) HandleCreateViewInteraction(ctx *fiber.Ctx) error {
	var params types.ViewQuestionParams

	if err := ctx.BodyParser(&params); err != nil {
		return err
	}

	interaction, err := h.interactionStore.CreateViewInteraction(ctx.Context(), &params)
	if err != nil {
		return err
	}

	return ctx.JSON(interaction)
}