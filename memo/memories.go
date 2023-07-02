package memo

import (
	"fmt"

	"github.com/google/uuid"
	pb "github.com/qdrant/go-client/qdrant"
	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type Memories struct {
	mongo  *mongo.Collection
	qdrant pb.PointsClient
	openai *openai.Client

	ctx context.Context
}

func (ms *Memories) AddMemory(id primitive.ObjectID, memory *Memory) (primitive.ObjectID, error) {
	res, err := ms.AddMemories(id, []*Memory{memory})
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res[0], nil
}

func (ms *Memories) AddMemories(id primitive.ObjectID, memories []*Memory) ([]primitive.ObjectID, error) {
	l := len(memories)
	var docs []interface{} = make([]interface{}, l)
	var mids []primitive.ObjectID = make([]primitive.ObjectID, l) // memory objectids
	var pids []uuid.UUID = make([]uuid.UUID, l)                   // point uuids
	var contents []string = make([]string, l)

	for idx, m := range memories {
		docs[idx] = m
		if m.ID == primitive.NilObjectID {
			m.ID = primitive.NewObjectID()
		}
		mids[idx] = m.ID
		contents[idx] = m.Content
		pid, _ := uuid.NewUUID()
		pids[idx] = pid
	}
	res, err := ms.mongo.InsertMany(ms.ctx, docs)
	if err != nil {
		return nil, err
	}

	if len(res.InsertedIDs) != l {
		return nil, fmt.Errorf("some memories not inserted: \n%v\n%v", res.InsertedIDs, mids)
	}

	// create embeddings
	ems, err := ms.embedding(contents)
	if err != nil {
		return nil, err
	}

	// upsert points into qdrant
	err = ms.upsertPoints(id, pids, mids, ems)
	return mids, nil
}

func (ms *Memories) DeleteMemory(id primitive.ObjectID) error {
	return nil
}

func (ms *Memories) DeleteMemories(ids []primitive.ObjectID) error {
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

func (ag *Memories) embedding(contents []string) ([]openai.Embedding, error) {
	// TODO: check the token limit
	// create embeddings
	req := openai.EmbeddingRequest{
		Input: contents,
		Model: openai.AdaEmbeddingV2,
	}

	res, err := ag.openai.CreateEmbeddings(ag.ctx, req)
	if err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (ag *Memories) upsertPoints(id primitive.ObjectID, pids []uuid.UUID, mids []primitive.ObjectID, ems []openai.Embedding) error {
	l := len(ems)
	points := make([]*pb.PointStruct, l)

	for _, em := range ems {
		point := &pb.PointStruct{
			Id:      &pb.PointId{PointIdOptions: &pb.PointId_Uuid{Uuid: pids[em.Index].String()}},
			Payload: map[string]*pb.Value{"mongo": {Kind: &pb.Value_StringValue{StringValue: mids[em.Index].Hex()}}},
			Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: em.Embedding}}},
		}
		points[em.Index] = point
	}
	waitUpsert := true
	_, err := ag.qdrant.Upsert(ag.ctx, &pb.UpsertPoints{
		CollectionName: id.Hex(),
		Wait:           &waitUpsert,
		Points:         points,
	})

	return err
}
