package memo

import (
	"fmt"
	"time"

	pb "github.com/qdrant/go-client/qdrant"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
)

// Agents is  a model which implements AgentModel interface
// it holds mongo collection and qdrant collection
type Agents struct {
	mongo  *mongo.Collection
	qdrant pb.CollectionsClient

	ListLimit int64
}

// Add agent and return inserted id
// if agent's id is set, it will return an error
// if agent's created time is not set, then it will be created
func (s *Agents) Add(ctx context.Context, agent *Agent) (primitive.ObjectID, error) {
	if agent.ID != primitive.NilObjectID {
		return primitive.NilObjectID, NewWrapError(400, fmt.Errorf("agent id should be nil"), "")
	}

	agent.ID = primitive.NewObjectID()
	if agent.Created.IsZero() {
		agent.Created = time.Now()
	}

	_, err := s.mongo.InsertOne(ctx, agent)
	if err != nil {
		return primitive.NilObjectID, err
	}

	err = s.createQdrantCollection(ctx, agent.ID.Hex())
	if err != nil {
		return primitive.NilObjectID, err
	}
	return agent.ID, nil
}

// Delete agent, if no agent matched it will return an notfound error
// id is agent's id
func (s Agents) Delete(ctx context.Context, id primitive.ObjectID) error {
	res, err := s.mongo.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	// check if agent exists
	if res.DeletedCount == 0 {
		return NewWrapError(404, fmt.Errorf("agent not found: %s", id), "")
	}

	return s.deleteQdrantCollection(ctx, id.Hex())
}

// Update an agent, if no agent matched it will return an notfound error
func (s *Agents) Update(ctx context.Context, agent *Agent) error {
	res, err := s.mongo.UpdateByID(ctx, agent.ID, bson.M{"$set": agent})
	if err != nil {
		return err
	}
	// if no agent matched
	if res.MatchedCount == 0 {
		return NewWrapError(404, fmt.Errorf("agent not found: %s", agent.ID.Hex()), "")
	}

	return err
}

// Get agent by id
func (s *Agents) Get(ctx context.Context, id primitive.ObjectID) (agent *Agent, err error) {
	err = s.mongo.FindOne(ctx, bson.M{"_id": id}).Decode(agent)
	if err == mongo.ErrNoDocuments {
		return nil, NewWrapError(404, fmt.Errorf("agent not found: %s", agent.ID.Hex()), "")
	}
	return
}

// List agents with offset, you can set search limit by session
func (s *Agents) List(ctx context.Context, offset primitive.ObjectID) (agents []*Agent, err error) {
	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(s.ListLimit)
	var filter bson.M
	// if offset is not nil, then make the offset filter
	if offset != primitive.NilObjectID {
		filter = bson.M{"_id": bson.M{"$lt": offset}}
	}
	cur, err := s.mongo.Find(ctx, filter, opts)
	if err != nil {
		return
	}

	err = cur.All(ctx, &agents)
	return
}

func (s Agents) createQdrantCollection(ctx context.Context, name string) (err error) {
	_, err = s.qdrant.Create(ctx, &pb.CreateCollection{
		CollectionName: name,
		VectorsConfig: &pb.VectorsConfig{
			Config: &pb.VectorsConfig_Params{
				Params: &pb.VectorParams{
					Size:     1536,
					Distance: pb.Distance_Cosine,
				},
			},
		},
	})
	return
}

func (s Agents) deleteQdrantCollection(ctx context.Context, name string) (err error) {
	_, err = s.qdrant.Delete(ctx, &pb.DeleteCollection{CollectionName: name})
	return
}
