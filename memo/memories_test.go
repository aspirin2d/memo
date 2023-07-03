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
	_, err = toml.DecodeFile("../.config.toml", &config)
	if err != nil {
		panic(err)
	}
	ms.memories = &Memories{
		qdrant: pb.NewPointsClient(qc),
		mongo:  mc.Database("test-db").Collection("memories"),
		openai: openai.NewClient(config.OpenAIAPIKey),

		Limit: 3, // search limit
	}
}

// create an agent before each test
func (ms *MemoriesSuite) SetupTest() {
	ctx := context.TODO()
	agent := &Agent{Name: "Aspirin"}
	ms.agents.Add(ctx, agent)
	ms.agent = agent
}

// delete the agent after each test
func (ms *MemoriesSuite) TearDownTest() {
	ctx := context.TODO()
	err := ms.agents.Delete(ctx, ms.agent.ID)
	ms.NoError(err)
}

func (ms *MemoriesSuite) TestAddMemory() {
	ctx := context.TODO()
	var memory1 = Memory{
		Content: "Hey, I am Aspirin",
	}

	id, err := ms.memories.AddOne(ctx, ms.agent.ID, &memory1)
	ms.NoError(err)

	var memory2 = Memory{
		Content: "Hey, I am Aspirin2D",
	}

	id, err = ms.memories.AddOne(ctx, ms.agent.ID, &memory2)
	ms.NoError(err)

	mem, err := ms.memories.GetOne(ctx, id)
	ms.NoError(err)

	pid := mem.PID
	aid := mem.AID

	// check the memory is added
	ms.Equal(id.Hex(), mem.ID.Hex())

	err = ms.memories.DeleteOne(ctx, id)
	ms.NoError(err)

	mem, err = ms.memories.GetOne(ctx, id)
	ms.ErrorIs(err, mongo.ErrNoDocuments)
	ms.Nil(mem)

	// check the memory's point is also delete
	res, err := ms.memories.qdrant.Get(ctx, &pb.GetPoints{
		CollectionName: aid.Hex(),
		Ids: []*pb.PointId{
			{PointIdOptions: &pb.PointId_Uuid{Uuid: pid}},
		}})
	ms.NoError(err)
	ms.Equal(0, len(res.Result))
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

	ids, err := ms.memories.AddMany(ctx, ms.agent.ID, memories)
	ms.NoError(err)

	// delete the first two memories
	err = ms.memories.DeleteMany(ctx, ids[:2])
	ms.NoError(err)

	agent, err := ms.memories.GetOne(ctx, ids[0])
	ms.ErrorIs(err, mongo.ErrNoDocuments)
	ms.Nil(agent)
}
func (ms *MemoriesSuite) TestSearchMemories() {
	ctx := context.TODO()
	var memories = []*Memory{
		{
			Content: "Hey, I am Aspirin.",
		},
		{
			Content: "My father is a teacher.",
		},
		{
			Content: "My favorite color is red.",
		},
		{
			Content: "My favorite food is pizza.",
		},
		{
			Content: "My favorite video game is Last of Us.",
		},
	}

	ids, err := ms.memories.AddMany(ctx, ms.agent.ID, memories)
	ms.NoError(err)
	ms.Equal(len(ids), len(memories))

	mems, err := ms.memories.Search(ctx, ms.agent.ID, "naughty dog")
	ms.NoError(err)
	ms.Equal(3, len(mems))

	ms.Contains(mems[0].Content, "video game")
}

func TestMemoriesSuite(t *testing.T) {
	suite.Run(t, new(MemoriesSuite))
}
