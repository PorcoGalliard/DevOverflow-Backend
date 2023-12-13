package db

const MongoDBName = "MONGO_DB_NAME"

type Store struct {
	Question QuestionStore
	User UserStore
	Tag TagStore
	Answer AnswerStore
	Interaction InteractionStore
}

type UserQueryParams struct {
	Page int64
	Limit int64
	Filter string
	SearchQuery string
}

type Map map[string]any