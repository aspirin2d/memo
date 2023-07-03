package memo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type vector []float32

type Memory struct {
	ID  primitive.ObjectID `bson:"_id" json:"id"`
	AID primitive.ObjectID `bson:"aid" json:"aid"` // agent's id
	PID string             `bson:"pid" json:"pid"` // memory's point id

	Content string    `bson:"content" json:"content"`
	Created time.Time `bson:"created_at" json:"created_at"`
}

type Agent struct {
	ID      primitive.ObjectID `bson:"_id" json:"id"`
	Name    string             `bson:"name" json:"name"`
	Created time.Time          `bson:"created_at" json:"created_at"`
}

type Config struct {
	OpenAIAPIKey string `toml:"openai_api_key"`
}
