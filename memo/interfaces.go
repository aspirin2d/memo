package memo

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

type AgentModel interface {
	// Add agent and return inserted id
	Add(ctx context.Context, agent *Agent) (primitive.ObjectID, error)

	// Delete agent by id
	Delete(ctx context.Context, id primitive.ObjectID) error

	// Update agent
	Update(ctx context.Context, agent *Agent) error

	// List and offset agent's id
	List(ctx context.Context, offset primitive.ObjectID) ([]*Agent, error)

	// Get agent by id
	Get(ctx context.Context, id primitive.ObjectID) (*Agent, error)
}

type MemoryModel interface {
	// Add memory and return inserted id
	AddOne(ctx context.Context, agent primitive.ObjectID, memory *Memory) (primitive.ObjectID, error)
	// Add memories and return inserted ids
	AddMany(ctx context.Context, agent primitive.ObjectID, memories []*Memory) ([]primitive.ObjectID, error)

	// Update memory
	UpdateOne(ctx context.Context, memory *Memory) error
	// Update memories
	// UpdateMany(ctx context.Context, memories []*Memory) error

	// Delete memory by id
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
	// Delete memories by ids
	DeleteMany(ctx context.Context, ids []primitive.ObjectID) error

	// ListMemories and offset memory's id
	// aid is agent's id which memories belong to
	// offset is the last memory's id
	List(ctx context.Context, aid primitive.ObjectID, offset primitive.ObjectID) ([]*Memory, error)

	// Search with query string, and return related memories and scores
	// aid is agent's id which memories belong to
	Search(ctx context.Context, aid primitive.ObjectID, query string) ([]*Memory, []float32, error)
}

// AgentController is a controller for handling agent requests
type AgentController interface {
	AddAgent(c *gin.Context)

	DeleteAgent(c *gin.Context)

	UpdateAgent(c *gin.Context)

	GetAgent(c *gin.Context)

	ListAgents(c *gin.Context)
}

// MemoryController is a controller for handling memory requests
type MemoryController interface {
	AddMemory(c *gin.Context)
	AddMemories(c *gin.Context)

	GetMemory(c *gin.Context)
	GetMemories(c *gin.Context)

	DeleteMemory(c *gin.Context)
	DeleteMemories(c *gin.Context)

	UpdateMemory(c *gin.Context)

	ListMemories(c *gin.Context)
}
