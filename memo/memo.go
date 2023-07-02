package memo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type vector []float32

type Memory struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Content string             `bson:"content" json:"content"`
	PID     string             `bson:"pid" json:"pid"`
	Created time.Time          `bson:"created_at" json:"created_at"`
}

type Agent struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Created time.Time          `bson:"created_at" json:"created_at"`
}

type Config struct {
	OpenAIAPIKey string `toml:"openai_api_key"`
}
