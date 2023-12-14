package api

import (
	"os"

	"github.com/gofiber/fiber/v2"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAIHandler struct {
	Client *openai.Client
}

func NewOpenAIHandler() *OpenAIHandler {
	return &OpenAIHandler{
		Client: openai.NewClient(os.Getenv("OPENAI_API_KEY")),
	}
}

func (h *OpenAIHandler) HandleChatGPT(ctx *fiber.Ctx) error {
	var reqBody map[string]string
	if err := ctx.BodyParser(&reqBody); err != nil {
		return err
	}

	question, ok := reqBody["question"]
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "missing question")
	}

	resp, err := h.Client.CreateChatCompletion(
		ctx.Context(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: "Hello dear my friend",
				},
				{
					Role: openai.ChatMessageRoleUser,
					Content: question,
				},
			},
		},
	)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	reply := resp.Choices[0].Message.Content
	return ctx.JSON(fiber.Map{"reply": reply})
}

