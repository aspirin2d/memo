package memo

import (
	"fmt"
	"testing"
	"time"

	pb "github.com/qdrant/go-client/qdrant"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AgentsSuite struct {
	suite.Suite
	agents *Agents
}

func (s *AgentsSuite) SetupSuite() {
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

	s.agents = &Agents{
		ctx:    ctx,
		qdrant: pb.NewCollectionsClient(qc),
		mongo:  mc.Database("test-db").Collection("agents"),
		limit:  15,
	}
}

func (s *AgentsSuite) TearDownTest() {
	// drop agent collection when each test finished
	err := s.agents.mongo.Drop(s.agents.ctx)
	s.NoError(err)
}
func (s *AgentsSuite) TearDownSuite() {
	// delete all qdrant collections when all tests in this suite finished
	res, err := s.agents.qdrant.List(s.agents.ctx, &pb.ListCollectionsRequest{})
	s.NoError(err)
	for _, col := range res.Collections {
		_, err := s.agents.qdrant.Delete(s.agents.ctx, &pb.DeleteCollection{CollectionName: col.Name})
		s.NoError(err)
	}
}

func (s *AgentsSuite) TestAddAgent() {
	id0, err := s.agents.AddAgent(&Agent{Name: "aspirin"})
	s.NoError(err)
	id1, err := s.agents.AddAgent(&Agent{Name: "aspirin2d"})
	s.NoError(err)

	s.NotEqual(id0.Hex(), id1.Hex())

	newID := primitive.NewObjectID()
	id3, err := s.agents.AddAgent(&Agent{ID: newID, Name: "aspirin2d"})
	s.NoError(err)
	s.Equal(newID.Hex(), id3.Hex())
}

func (s *AgentsSuite) TestGetAgent() {
	id, err := s.agents.AddAgent(&Agent{Name: "aspirin"})
	s.NoError(err)
	agent, err := s.agents.GetAgent(id)
	s.NoError(err)
	s.Equal(agent.Name, "aspirin")
	// it will create a new "created" value for the agent
	s.True(agent.Created.After(time.Now().Add(-5 * time.Second)))

	_, err = s.agents.GetAgent(primitive.NewObjectID())
	s.Error(err)
}

func (s *AgentsSuite) TestListAgents() {
	for i := range [5]int{} {
		_, err := s.agents.AddAgent(&Agent{Name: fmt.Sprintf("aspirin %d", i)})
		s.NoError(err)
	}
	agents, err := s.agents.ListAgents(primitive.NilObjectID)
	s.NoError(err)
	s.Equal(5, len(agents))

	for i := range [20]int{} {
		_, err := s.agents.AddAgent(&Agent{Name: fmt.Sprintf("aspirin %d", i)})
		s.NoError(err)
	}

	agents, err = s.agents.ListAgents(primitive.NilObjectID)
	s.NoError(err)
	// reached search limit
	s.Equal(15, len(agents))

	// search with the last agent as offset, it will get the rest of the agents
	agents, err = s.agents.ListAgents(agents[len(agents)-1].ID)
	s.NoError(err)
	s.Equal(10, len(agents))
}

func (s *AgentsSuite) TestDeleteAgent() {
	for i := range [5]int{} {
		_, err := s.agents.AddAgent(&Agent{Name: fmt.Sprintf("aspirin %d", i)})
		s.NoError(err)
	}

	agents, err := s.agents.ListAgents(primitive.NilObjectID)
	s.NoError(err)
	s.Equal(5, len(agents))

	s.agents.DeleteAgent(agents[len(agents)-1].ID)
	agents, err = s.agents.ListAgents(primitive.NilObjectID)
	s.NoError(err)
	s.Equal(4, len(agents))

	err = s.agents.DeleteAgent(primitive.NewObjectID())
	s.Error(err)
}

func (s *AgentsSuite) TestUpdateAgent() {
	id, err := s.agents.AddAgent(&Agent{Name: "aspirin2d"})
	s.NoError(err)
	err = s.agents.UpdateAgent(&Agent{ID: id, Name: "aspirin3d"})
	s.NoError(err)
	agent, err := s.agents.GetAgent(id)
	s.NoError(err)
	s.Equal("aspirin3d", agent.Name)

	// try to update a not existed agent will cause error
	err = s.agents.UpdateAgent(&Agent{ID: primitive.NewObjectID(), Name: "aspirin3d"})
	s.Error(err)
}

func TestAgentsSuite(t *testing.T) {
	suite.Run(t, new(AgentsSuite))
}
