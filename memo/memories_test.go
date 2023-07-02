package memo

import (
	"testing"

	"github.com/BurntSushi/toml"
	pb "github.com/qdrant/go-client/qdrant"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MemoriesSuite struct {
	suite.Suite
	agents   *Agents
	memories *Memories

	agent *Agent
}

func (ms *MemoriesSuite) SetupSuite() {
	// check agents is implements AgentModel
	var _ AgentModel = (*Agents)(nil)

	ctx := context.TODO()
	// mongodb
	mc, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	// qdrant
	qc, err := grpc.Dial("localhost:6334", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	ms.agents = &Agents{
		qdrant: pb.NewCollectionsClient(qc),
		mongo:  mc.Database("test-db").Collection("agents"),
		limit:  15,
	}
	var config Config
	_, err = toml.DecodeFile("../config.toml", &config)
	if err != nil {
		panic(err)
	}
	ms.memories = &Memories{
		qdrant: pb.NewPointsClient(qc),
		mongo:  mc.Database("test-db").Collection("memories"),
		openai: openai.NewClient(config.OpenAIAPIKey),
	}
}

// create an agent before each test
func (ms *MemoriesSuite) SetupTest() {
	ctx := context.TODO()
	agent := &Agent{Name: "Aspirin"}
	ms.agents.AddAgent(ctx, agent)
	ms.agent = agent
}

// delete the agent after each test
func (ms *MemoriesSuite) TearDownTest() {
	ctx := context.TODO()
	err := ms.agents.DeleteAgent(ctx, ms.agent.ID)
	ms.NoError(err)
}

func (ms *MemoriesSuite) TestAddMemory() {
	ctx := context.TODO()
	var memory = Memory{
		Content: "Hey, I am Aspirin",
	}

	id, err := ms.memories.AddMemory(ctx, ms.agent.ID, &memory)
	ms.NoError(err)

	ag, err := ms.memories.GetMemory(ctx, id)
	ms.NoError(err)

	// check the memory is added
	ms.Equal(id.Hex(), ag.ID.Hex())

	err = ms.memories.DeleteMemory(ctx, id)
	ms.NoError(err)

	agent, err := ms.memories.GetMemory(ctx, id)
	ms.ErrorIs(err, mongo.ErrNoDocuments)
	ms.Nil(agent)
}

func (ms *MemoriesSuite) TestAddMemories() {
	ctx := context.TODO()
	var memories = []*Memory{
		{
			Content: "Hey, I am Aspirin",
		},
		{
			Content: "My father is a teacher",
		},
		{
			Content: "My favorite color is red",
		},
		{
			Content: "My favorite food is pizza",
		},
		{
			Content: "I like to play video games",
		},
	}

	ids, err := ms.memories.AddMemories(ctx, ms.agent.ID, memories)
	ms.NoError(err)

	// delete the first two memories
	err = ms.memories.DeleteMemories(ctx, ids[:2])
	ms.NoError(err)

	agent, err := ms.memories.GetMemory(ctx, ids[0])
	ms.ErrorIs(err, mongo.ErrNoDocuments)
	ms.Nil(agent)
}

func TestMemoriesSuite(t *testing.T) {
	suite.Run(t, new(MemoriesSuite))
}
