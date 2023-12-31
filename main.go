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
	"github.com/sashabaranov/go-openai"
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
	openAIAPIKey := os.Getenv("OPENAI_API_KEY")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	} 
	openAIClient := openai.NewClient(openAIAPIKey)

	var (
		userStore = db.NewMongoUserStore(client)
		tagStore = db.NewMongoTagStore(client)
		answerStore = db.NewMongoAnswerStore(client, userStore)
		questionStore = db.NewMongoQuestionStore(client, tagStore, userStore)
		interactionStore = db.NewMongoInteractionStore(client, userStore, questionStore)

		store = &db.Store{
			Question: questionStore,
			User: userStore,
			Tag: tagStore,
			Answer: answerStore,
			Interaction: interactionStore,
		}

		openAIHandler = api.NewOpenAIHandler(openAIClient)
		questionHandler = api.NewQuestionHandler(store.Question, store.User, store.Tag, store.Answer)
		userHandler = api.NewUserHandler(store.User, store.Tag, store.Question)
		tagHandler = api.NewTagHandler(store.Tag, store.User)
		answerHandler = api.NewAnswerHandler(store.Answer, store.Question, store.User)
		interactionHandler = api.NewInteractionHandler(store.Interaction)
		app = fiber.New(config)
		auth = app.Group("/api")
		apiv1 = app.Group("/api/v1")
	)

	app.Use(cors.New())
	// Question Handler
	apiv1.Get("/question/:id", questionHandler.HandleGetQuestionByID)
	apiv1.Get("/question", questionHandler.HandleGetQuestions)
	apiv1.Get("/question/user/:id", questionHandler.HandleGetQuestionsByUserID)
	apiv1.Post("/ask-question", questionHandler.HandleAskQuestion)
	apiv1.Post("/question/:id/vote", questionHandler.HandleQuestionVote)
	apiv1.Delete("/question/:_id", questionHandler.HandleDeleteQuestionByID)
	
	// User Handler
	app.Get("/", userHandler.HandleSayHello)
	apiv1.Get("/user/:clerkID", userHandler.HandleGetUserByID)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:clerkID/saved-questions", questionHandler.HandleGetSavedQuestions)
	auth.Post("/sign-up", userHandler.HandleCreateUser)
	apiv1.Post("/user/save-question", userHandler.HandleSaveQuestion)
	apiv1.Put("/user/:clerkID", userHandler.HandleUpdateUser)
	apiv1.Delete("/user/:clerkID", userHandler.HandleDeleteUser)

	// Tag Handler
	apiv1.Get("/tag/:_id", tagHandler.HandleGetTagByID)
	apiv1.Get("/tag/:name", tagHandler.HandleGetTagByName)
	apiv1.Get("/tag", tagHandler.HandleGetTags)
	apiv1.Get("/tag/:id/questions", questionHandler.HandleGetQuestiosByTagID)
	apiv1.Post("/tag", tagHandler.HandleCreateTag)
	apiv1.Put("/tag/:_id", tagHandler.HandleUpdateTag)

	// Answer Handler
	apiv1.Get("/question/:questionID/answer/:answerID", answerHandler.HandleGetAnswerByID)
	apiv1.Get("/question/:id/answers", answerHandler.HandleGetAnswersOfQuestion)
	apiv1.Get("/answer/user/:id", answerHandler.HandleGetAnswersByUserID)
	apiv1.Post("/answer/:id/vote", answerHandler.HandleAnswerVote)
	apiv1.Post("/answer-question", answerHandler.HandleCreateAnswer)

	// Interaction Handler
	apiv1.Post("/question/view", interactionHandler.HandleCreateViewInteraction)

	// OpenAI Handler
	apiv1.Post("/chat-gpt", openAIHandler.HandleChatGPT)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000" // Default port if not specified
	}
	listenAddr := ":" + port

	app.Listen(listenAddr)
}

func init() {
	_, isHeroku := os.LookupEnv("DYNO")
	if !isHeroku {
		if err := godotenv.Load(); err != nil {
			log.Fatal(err)
		}
	}
}