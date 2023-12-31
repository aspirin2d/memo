package memo

import (
	"testing"

	"github.com/BurntSushi/toml"
	pb "github.com/qdrant/go-client/qdrant"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	var _ MemoryModel = (*Memories)(nil)

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
		qdrant:    pb.NewCollectionsClient(qc),
		mongo:     mc.Database("test-db").Collection("agents"),
		ListLimit: 15,
	}
	var config Config
	_, err = toml.DecodeFile("../.config.toml", &config)
	if err != nil {
		panic(err)
	}
	ms.memories = &Memories{
		qdrant:      pb.NewPointsClient(qc),
		mongo:       mc.Database("test-db").Collection("memories"),
		llm:         NewOpenAI(config.OpenAIAPIKey),
		SearchLimit: 3, // search limit
	}
}

// create an agent before each test
func (ms *MemoriesSuite) SetupTest() {
	ctx := context.TODO()
	agent := &Agent{Name: "Aspirin"}
	_, err := ms.agents.Add(ctx, agent)
	ms.NoError(err)
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

	_, err := ms.memories.AddOne(ctx, ms.agent.ID, &memory1)
	ms.NoError(err)

	var memory2 = Memory{
		Content: "Hey, I am Aspirin2D",
	}

	id, err := ms.memories.AddOne(ctx, ms.agent.ID, &memory2)
	ms.NoError(err)

	mem, err := ms.memories.GetOne(ctx, ms.agent.ID, id)
	ms.NoError(err)

	pid := mem.PID
	aid := mem.AID

	// check the memory is added
	ms.Equal(id.Hex(), mem.ID.Hex())

	err = ms.memories.DeleteOne(ctx, ms.agent.ID, id)
	ms.NoError(err)

	_, err = ms.memories.GetOne(ctx, ms.agent.ID, id)
	ms.Error(err, mongo.ErrNoDocuments)

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
	err = ms.memories.DeleteMany(ctx, ms.agent.ID, ids[:2])
	ms.NoError(err)

	_, err = ms.memories.GetOne(ctx, ms.agent.ID, ids[0])
	ms.Error(err)

	mems, err := ms.memories.GetMany(ctx, ms.agent.ID, ids[2:])
	ms.NoError(err)
	ms.Len(mems, 3)
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

	mems, scores, err := ms.memories.Search(ctx, ms.agent.ID, "naughty dog") // just for fun
	ms.NoError(err)
	ms.Len(mems, 3)
	ms.Len(scores, 3)

	ms.Contains(mems[0].Content, "Last of Us")
}

func (ms *MemoriesSuite) TestUpdateMemory() {
	ctx := context.TODO()
	var memory1 = Memory{
		Content: "Hey, I am Aspirin",
	}

	id, err := ms.memories.AddOne(ctx, ms.agent.ID, &memory1)
	ms.NoError(err)

	memory1.Content = "Hey, I am Aspirin2D"

	err = ms.memories.UpdateOne(ctx, ms.agent.ID, &memory1)
	ms.NoError(err)

	mem, err := ms.memories.GetOne(ctx, ms.agent.ID, id)
	ms.NoError(err)

	ms.Equal("Hey, I am Aspirin2D", mem.Content)
}

func (ms *MemoriesSuite) TestUpdateMemory2() {
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

	ids, _ := ms.memories.AddMany(ctx, ms.agent.ID, memories)
	memories, _ = ms.memories.GetMany(ctx, ms.agent.ID, ids)
	memories[0].Content = "The sky is blue"
	memories[1].Content = "The water is green"
	memories[2].Content = "The moon is orange"
	memories[3].Content = "The apple is red"
	memories[4].Content = "The sand is yellow"
	_ = ms.memories.UpdateMany(ctx, ms.agent.ID, memories)

	mems, _, _ := ms.memories.Search(ctx, ms.agent.ID, "planet") // just for fun
	ms.Contains(mems[0].Content, "moon")
}

func (ms *MemoriesSuite) TestListMemories() {
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

	ms.memories.ListLimit = 3 // set 3 per page
	mems, err := ms.memories.List(ctx, ms.agent.ID, primitive.NilObjectID)
	ms.NoError(err)
	ms.Len(mems, 3)

	mems, err = ms.memories.List(ctx, ms.agent.ID, mems[2].ID)
	ms.NoError(err)
	ms.Len(mems, 2)
}

func TestMemoriesSuite(t *testing.T) {
	suite.Run(t, new(MemoriesSuite))
}
