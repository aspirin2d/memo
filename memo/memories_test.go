package memo

import (
	"testing"

	pb "github.com/qdrant/go-client/qdrant"
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
		ctx:    ctx,
		qdrant: pb.NewCollectionsClient(qc),
		mongo:  mc.Database("test-db").Collection("agents"),
		limit:  15,
	}

	ms.memories = &Memories{
		ctx: ctx,
	}
}

func TestMemoriesSuite(t *testing.T) {
	suite.Run(t, new(MemoriesSuite))
}
