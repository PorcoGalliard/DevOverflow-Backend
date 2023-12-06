package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fullstack/dev-overflow/db"
	"github.com/fullstack/dev-overflow/db/fixtures"

	// "github.com/fullstack/dev-overflow/types"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	var (
		mongoDBEndpoint = os.Getenv("MONGO_DB_URL")
		mongoDBName = os.Getenv("MONGO_DB_NAME")
	)

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDBEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	if err := mongoClient.Database(mongoDBName).Drop(context.Background()); err != nil {
		log.Fatal(err)
	}

	userStore := db.NewMongoUserStore(mongoClient)
	tagStore := db.NewMongoTagStore(mongoClient)
	questionStore := db.NewMongoQuestionStore(mongoClient, tagStore, userStore)
	store := &db.Store{
		Question: questionStore,
		User: userStore,
		Tag: tagStore,
	}


	user := fixtures.AddUser(store, "Higuruma", "Hiromi", "12345678", "higuruma@gmail.com", "higurumahiromu123")
	fmt.Println("User berhasil ditambahkan, berikut adalah ID user =>", user.ID)

	tag := fixtures.AddTag(store, "Go")
	fmt.Println("Tag berhasil ditambahkan, berikut adalah ID tag =>", tag.ID)

	question := fixtures.AddQuestion(store, "Bagaimana cara menggunakan Go?", "Saya baru belajar Go dan ingin tahu lebih banyak.", user.ID, []primitive.ObjectID{tag.ID}, time.Now())
	fmt.Println("Pertanyaan berhasil ditambahkan, berikut adalah ID pertanyaan =>", question.ID)
}