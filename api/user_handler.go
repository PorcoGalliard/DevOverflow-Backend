package api

import (
	"errors"
	"os"
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/fullstack/dev-overflow/db"
	"github.com/fullstack/dev-overflow/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
	tagStore db.TagStore
	questionStore db.QuestionStore
}

func NewUserHandler(userStore db.UserStore, tagStore db.TagStore, questionStore db.QuestionStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		tagStore: tagStore,
		questionStore: questionStore,
	}
}

func (h *UserHandler) HandleSignUp(c *fiber.Ctx) error {
	sessionToken := c.Get("Authorization")
	apiKey := os.Getenv("CLERK_SECRET_KEY")

	clerkClient, err := clerk.NewClient(apiKey)
	if err != nil {
		return ErrBadRequest()
	}

	session, err := clerkClient.Sessions().Verify(sessionToken, "")
	if err != nil {
		return ErrUnauthorized()
	}
	
	clerkUser, err := clerkClient.Users().Read(session.ID)
	if err != nil {
		return ErrInvalidID()
	}

	userOid := primitive.NewObjectID()

	user := &types.User{
		ID: userOid,
		ClerkID: clerkUser.ID,
		JoinedAt: time.Now().UTC(),
	}

	insertedUser, err := h.userStore.CreateUser(c.Context(), user)
	if err != nil {
		return ErrBadRequest()
	}

	return c.JSON(fiber.Map{"message": "User berhasil ditambahkan, ID nya adalah =>" + insertedUser.ClerkID})
}

func (h *UserHandler) HandleGetUserByID(ctx *fiber.Ctx) error {
	var (
		id = ctx.Params("clerkID")
	) 

	user, err := h.userStore.GetUserByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound(id)
		}
	}

	return ctx.JSON(user)
}

func (h *UserHandler) HandleGetUsers(ctx *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(ctx.Context())
	if err != nil {
		return ErrBadRequest()
	}

	return ctx.JSON(users)
}

func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var (
		params types.CreateUserParam
	)

	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return ErrBadRequest()
	}

	insertedUser, err := h.userStore.CreateUser(c.Context(), user)
	if err != nil {
		return ErrBadRequest()
	}

	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	var (
		clerkID = c.Params("clerkID")
		params *types.UpdateUserParam
	)

	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	if params.UpdateData["email"] != nil || params.UpdateData["firstName"] != nil || params.UpdateData["lastName"] != nil {
		if errors := params.Validate(); len(errors) > 0 {
			return c.JSON(errors)
		}
	}

	if err := h.userStore.UpdateUser(c.Context(), clerkID, params); err != nil {
		return ErrBadRequest()
	}

	return c.JSON(map[string]string{"message": "User berhasil diupdate dengan ID => " + clerkID})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	var (
		clerkID = c.Params("clerkID")
	)

	user, err := h.userStore.GetUserByID(c.Context(), clerkID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrResourceNotFound(clerkID)
		}
	}

	if err := h.tagStore.UpdateManyFollowersByID(c.Context(), user.ID); err != nil {
		return ErrBadRequest()
	}

	questions, err := h.questionStore.GetQuestionsByUserID(c.Context(), clerkID)
	if err != nil {
		return ErrBadRequest()
	}

	for i := 0; i < len(questions); i++ {
		if err := h.tagStore.UpdateManyQuestionsByID(c.Context(), questions[i].ID); err != nil {
			return ErrBadRequest()
		}
	}

	if err := h.questionStore.DeleteManyQuestionsByUserID(c.Context(), user.ID); err != nil {
		return ErrBadRequest()
	}

	if err := h.userStore.DeleteUser(c.Context(), clerkID); err != nil {
		return ErrBadRequest()
	}

	return c.JSON(map[string]string{"message": "User berhasil dihapus dengan ID => " + clerkID})
}

func (h *UserHandler) HandleSayHello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello World"})
}