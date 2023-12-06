package main

import (
	"context"
	"log"
	"os"

	"github.com/fullstack/dev-overflow/api"
	"github.com/fullstack/dev-overflow/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	config = fiber.Config{
		ErrorHandler: api.ErrorHandler,
	}
)

func main() {
	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	} 

	var (
		userStore = db.NewMongoUserStore(client)
		tagStore = db.NewMongoTagStore(client)
		questionStore = db.NewMongoQuestionStore(client, tagStore, userStore)

		store = &db.Store{
			Question: questionStore,
			User: userStore,
			Tag: tagStore,
		}

		questionHandler = api.NewQuestionHandler(store.Question, store.User, store.Tag)
		userHandler = api.NewUserHandler(store.User, store.Tag, store.Question)
		tagHandler = api.NewTagHandler(store.Tag, store.User)
		app = fiber.New(config)
		auth = app.Group("/api")
		apiv1 = app.Group("/api/v1")
	)

	app.Use(cors.New())
	// Question Handler
	apiv1.Get("/question/:_id", questionHandler.HandleGetQuestionByID)
	apiv1.Get("/question", questionHandler.HandleGetQuestions)
	apiv1.Get("/question/user/:_id", questionHandler.HandleGetQuestionsByUserID)
	apiv1.Post("/ask-question", questionHandler.HandleAskQuestion)
	apiv1.Delete("/question/:_id", questionHandler.HandleDeleteQuestionByID)
	

	// User Handler
	// auth.Post("/sign-up", userHandler.HandleSignUp)
	apiv1.Get("/user/:clerkID", userHandler.HandleGetUserByID)
	auth.Post("/sign-up", userHandler.HandleCreateUser)
	apiv1.Put("/user/:clerkID", userHandler.HandleUpdateUser)
	apiv1.Delete("/user/:clerkID", userHandler.HandleDeleteUser)

	// Tag Handler
	apiv1.Get("/tag/:_id", tagHandler.HandleGetTagByID)
	apiv1.Get("/tag/:name", tagHandler.HandleGetTagByName)
	apiv1.Get("/tag", tagHandler.HandleGetTags)
	apiv1.Post("/tag", tagHandler.HandleCreateTag)
	apiv1.Put("/tag/:_id", tagHandler.HandleUpdateTag)
	

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(listenAddr)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}