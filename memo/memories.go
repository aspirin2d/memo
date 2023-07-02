package memo

import (
	"context"
	"fmt"

	pb "github.com/qdrant/go-client/qdrant"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Memories is a model which implements MemoryModel interface
// mongo is a mongo collection of memories
// qdrant is a qdrant points of memories
// openai is an openai client
type Memories struct {
	mongo  *mongo.Collection
	qdrant pb.PointsClient
	openai *openai.Client
}

// AddMemory adds a memory to the agent
// aid is agent's id
func (ms *Memories) AddMemory(ctx context.Context, aid primitive.ObjectID, memory *Memory) (primitive.ObjectID, error) {
	res, err := ms.AddMemories(ctx, aid, []*Memory{memory})
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res[0], nil
}

// AddMemories adds memories to the agent
// aid is agent's id
func (ms *Memories) AddMemories(ctx context.Context, id primitive.ObjectID, memories []*Memory) ([]primitive.ObjectID, error) {
	l := len(memories)
	var docs []interface{} = make([]interface{}, l)
	var mids []primitive.ObjectID = make([]primitive.ObjectID, l) // memory objectids
	var contents []string = make([]string, l)

	for idx, m := range memories {
		docs[idx] = m
		if m.ID != primitive.NilObjectID {
			return nil, fmt.Errorf("memory id should be nil")
		}
		m.ID = primitive.NewObjectID()
		mids[idx] = m.ID
		contents[idx] = m.Content
	}

	res, err := ms.mongo.InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}

	if len(res.InsertedIDs) != l {
		return nil, fmt.Errorf("some memories not inserted: \n%v\n%v", res.InsertedIDs, mids)
	}

	// create embeddings
	ems, err := ms.embedding(ctx, contents)
	if err != nil {
		return nil, err
	}

	// upsert points into qdrant
	err = ms.upsertPoints(ctx, id, mids, ems)
	return mids, nil
}

// GetMemory gets a memory by id
func (ms *Memories) GetMemory(ctx context.Context, id primitive.ObjectID) (memory *Memory, err error) {
	err = ms.mongo.FindOne(ctx, bson.M{"_id": id}).Decode(&memory)
	if err != nil {
		return nil, err
	}
	return memory, nil
}

// GetMemories gets memories by ids
func (ms *Memories) GetMemories(ctx context.Context, ids []primitive.ObjectID) (memories []*Memory, err error) {
	cur, err := ms.mongo.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	err = cur.All(ctx, &memories)
	if err != nil {
		return nil, err
	}
	if len(memories) != len(ids) {
		return nil, fmt.Errorf("some memories not found: \n%v", ids)
	}
	return
}

// DeleteMemory deletes a memory by id
func (ms *Memories) DeleteMemory(ctx context.Context, id primitive.ObjectID) error {
	_, err := ms.mongo.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

// DeleteMemories deletes memories by ids
func (ms *Memories) DeleteMemories(ctx context.Context, ids []primitive.ObjectID) error {
	res, err := ms.mongo.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return err
	}
	if res.DeletedCount != int64(len(ids)) {
		return fmt.Errorf("some memories not deleted: \n%v", ids)
	}
	return nil
}

// func (ag *Agent) UpdateMemory(memory Memory) error {

// }
// func (ag *Agent) UpdateMemories(memories []Memory) error {

// }

// func (ag *Agent) ListMemories(offset primitive.ObjectID) ([]Memory, error) {

// }

// func (ag *Agent) Search(query string, limit string) ([]Memory, error) {

// }

func (ag *Memories) embedding(ctx context.Context, contents []string) ([]openai.Embedding, error) {
	// TODO: check the token limit
	// create embeddings
	req := openai.EmbeddingRequest{
		Input: contents,
		Model: openai.AdaEmbeddingV2,
	}

	res, err := ag.openai.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (ag *Memories) upsertPoints(ctx context.Context, id primitive.ObjectID, mids []primitive.ObjectID, ems []openai.Embedding) error {
	l := len(ems)
	points := make([]*pb.PointStruct, l)
	for _, em := range ems {
		point := &pb.PointStruct{
			Id:      &pb.PointId{PointIdOptions: &pb.PointId_Num{Num: uint64(mids[em.Index].Timestamp().Unix())}},
			Payload: map[string]*pb.Value{"mongo": {Kind: &pb.Value_StringValue{StringValue: mids[em.Index].Hex()}}},
			Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: em.Embedding}}},
		}
		points[em.Index] = point
	}
	waitUpsert := true
	_, err := ag.qdrant.Upsert(ctx, &pb.UpsertPoints{
		CollectionName: id.Hex(),
		Wait:           &waitUpsert,
		Points:         points,
	})

	return err
}
