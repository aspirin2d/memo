package memo

import (
	"context"
	"time"

	"github.com/BurntSushi/toml"
	pb "github.com/qdrant/go-client/qdrant"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const AGENTS_COLLECTION = "agents"
const MEMORIES_COLLECTION = "memories"

type Memory struct {
	ID  primitive.ObjectID `bson:"_id" json:"id"`
	AID primitive.ObjectID `bson:"aid" json:"aid"` // agent's id
	PID string             `bson:"pid" json:"pid"` // memory's point id

	Content string    `bson:"content" json:"content"`
	Created time.Time `bson:"created_at" json:"created_at"`
}

type Agent struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name    string             `bson:"name" json:"name"`
	Created time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}

type Config struct {
	OpenAIAPIKey string `toml:"openai_api_key"`

	MongoUri  string `toml:"mongo_uri"`
	MongoDb   string `toml:"mongo_db"`
	QdrantUri string `toml:"qdrant_uri"`

	AgentListLimit    int `toml:"agent_search_limit"`
	MemorySearchLimit int `toml:"memory_search_limit"`
	MemoryListLimit   int `toml:"memory_list_limit"`
}

type Memo struct {
	Config *Config

	Agents   AgentModel  // agents model
	Memories MemoryModel // memories model

	Logger *zap.SugaredLogger
}

// NewAgents creates the default Agents which implements the AgentModel interface.
// config_path is the path to the config file.
func FromConfig(config_path string) *Memo {
	var conf Config = Config{
		MongoUri:          "mongodb://localhost:27017",
		MongoDb:           "memo",
		QdrantUri:         "localhost:6334",
		AgentListLimit:    15,
		MemoryListLimit:   15,
		MemorySearchLimit: 5, // top_k
	}

	_, err := toml.DecodeFile(config_path, &conf)
	if err != nil {
		panic(err)
	}

	if conf.OpenAIAPIKey == "" {
		panic("OpenAIAPIKey is empty")
	}

	ctx := context.TODO()
	// mongodb
	mc, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.MongoUri))
	if err != nil {
		panic(err)
	}

	// qdrant
	qc, err := grpc.Dial(conf.QdrantUri, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	// logger
	logger, _ := zap.NewProduction()
	defer func() {
		err = logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	return &Memo{
		Agents: &Agents{
			mongo:     mc.Database(conf.MongoDb).Collection(AGENTS_COLLECTION),
			qdrant:    pb.NewCollectionsClient(qc),
			ListLimit: int64(conf.AgentListLimit),
		},
		Memories: &Memories{
			mongo:       mc.Database(conf.MongoDb).Collection(MEMORIES_COLLECTION),
			qdrant:      pb.NewPointsClient(qc),
			SearchLimit: int64(conf.MemorySearchLimit),
			ListLimit:   int64(conf.MemoryListLimit),
		},

		Logger: logger.Sugar(),
	}
}
