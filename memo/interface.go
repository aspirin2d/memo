package memo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

type AgentModel interface {
	// Add agent and return inserted id
	AddAgent(ctx context.Context, agent *Agent) (primitive.ObjectID, error)

	// Delete agent by id
	DeleteAgent(ctx context.Context, id primitive.ObjectID) error

	// Update agent
	UpdateAgent(ctx context.Context, agent *Agent) error

	// ListAgents and offset agent's id
	ListAgents(ctx context.Context, offset primitive.ObjectID) ([]*Agent, error)

	// GetAgent by id
	GetAgent(ctx context.Context, id primitive.ObjectID) (*Agent, error)
}

type MemoryModel interface {
	// Add memory and return inserted id
	AddMemory(ctx context.Context, agent primitive.ObjectID, memory Memory) (string, error)
	// Add memories and return inserted ids
	AddMemories(ctx context.Context, agent primitive.ObjectID, memories []Memory) ([]string, error)

	// Update memory
	UpdateMemory(memory Memory) error
	// Update memories
	UpdateMemories(memories []Memory) error

	// Delete memory by id
	DeleteMemory(id primitive.ObjectID) error
	// Delete memories by ids
	DeleteMemories(ids []primitive.ObjectID) error

	// ListMemories and offset memory's id
	ListMemories(offset primitive.ObjectID) ([]Memory, string, error)

	// Search related memories by query string,
	Search(query string, limit string) ([]Memory, error)
}
