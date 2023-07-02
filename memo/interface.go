package memo

import (
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

	// Get by id
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
	UpdateMany(ctx context.Context, memories []*Memory) error

	// Delete memory by id
	DeleteOne(ctx context.Context, id primitive.ObjectID) error
	// Delete memories by ids
	DeleteMany(ids []primitive.ObjectID) error

	// ListMemories and offset memory's id
	// offset is the last memory's id
	List(ctx context.Context, offset primitive.ObjectID) ([]*Memory, error)

	// Search related memories by query string,
	Search(ctx context.Context, query string, limit string) ([]*Memory, error)
}
