package memo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

// AddOne adds a memory to the agent
// aid is agent's id
func (ms *Memories) AddOne(ctx context.Context, aid primitive.ObjectID, memory *Memory) (primitive.ObjectID, error) {
	res, err := ms.AddMany(ctx, aid, []*Memory{memory})
	if err != nil {
		return primitive.NilObjectID, err
	}
	return res[0], nil
}

// AddMany adds memories to the agent
// aid is agent's id
func (ms *Memories) AddMany(ctx context.Context, aid primitive.ObjectID, memories []*Memory) ([]primitive.ObjectID, error) {
	l := len(memories)

	var docs []interface{} = make([]interface{}, l)               // mongodb documents
	var contents []string = make([]string, l)                     // memory contents
	var mids []primitive.ObjectID = make([]primitive.ObjectID, l) // memory objectids
	var pids []uuid.UUID = make([]uuid.UUID, l)                   // point ids

	for idx, m := range memories {
		docs[idx] = m
		// check if memory id is nil
		if m.ID != primitive.NilObjectID {
			return nil, fmt.Errorf("memory id should be nil")
		}

		// check if memory aid is nil
		if m.AID != primitive.NilObjectID {
			return nil, fmt.Errorf("memory agent id should be nil")
		}

		m.ID = primitive.NewObjectID()
		m.AID = aid

		mids[idx] = m.ID
		contents[idx] = m.Content

		// create a reference to the point
		pids[idx] = uuid.New()
		m.PID = pids[idx].String()
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
	err = ms.upsertPoints(ctx, aid, pids, mids, ems)
	return mids, err
}

// GetOne gets a memory by id
func (ms *Memories) GetOne(ctx context.Context, id primitive.ObjectID) (memory *Memory, err error) {
	err = ms.mongo.FindOne(ctx, bson.M{"_id": id}).Decode(&memory)
	if err != nil {
		return nil, err
	}
	return memory, nil
}

// GetMany gets memories by ids
func (ms *Memories) GetMany(ctx context.Context, ids []primitive.ObjectID) (memories []*Memory, err error) {
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

// DeleteOne deletes a memory by id
func (ms *Memories) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	var mem Memory
	err := ms.mongo.FindOneAndDelete(ctx, bson.M{"_id": id}).Decode(&mem)
	if err != nil {
		return err
	}

	// also need to delete it from qdrant
	uid, err := uuid.Parse(mem.PID)
	if err != nil {
		return err
	}
	err = ms.deletePoints(ctx, mem.AID, []uuid.UUID{uid})
	return err
}

// DeleteMany deletes memories by ids
func (ms *Memories) DeleteMany(ctx context.Context, ids []primitive.ObjectID) error {
	res, err := ms.mongo.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return err
	}
	var mems []*Memory
	if err = res.All(ctx, &mems); err != nil {
		return err
	}

	// check if all memories found
	if len(mems) != len(ids) {
		return fmt.Errorf("some memories not found: \n%v", ids)
	}

	// delete memories from mongodb
	dres, err := ms.mongo.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if dres.DeletedCount != int64(len(ids)) {
		return fmt.Errorf("some memories not deleted: \n%v", ids)
	}

	// get point ids
	var pids []uuid.UUID = make([]uuid.UUID, len(mems))
	for idx, m := range mems {
		pids[idx], err = uuid.Parse(m.PID)
		if err != nil {
			return err
		}
	}

	// finally delete memories' points from qdrant
	return ms.deletePoints(ctx, mems[0].AID, pids)
}

// func (ag *Agent) UpdateMemory(memory Memory) error {

// }
// func (ag *Agent) UpdateMemories(memories []Memory) error {

// }

// func (ag *Agent) ListMemories(offset primitive.ObjectID) ([]Memory, error) {

// }

// Search searches memories by query
// id is aeget's id
func (ms *Memories) Search(ctx context.Context, id primitive.ObjectID, query string, limit string) ([]*Memory, error) {
	ems, err := ms.embedding(ctx, []string{query})
	if err != nil {
		return nil, err
	}
	res, err := ms.qdrant.Search(ctx, &pb.SearchPoints{
		CollectionName: id.Hex(),
		Vector:         ems[0].Embedding,
		WithPayload:    &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{Enable: true}},  // with payload
		WithVectors:    &pb.WithVectorsSelector{SelectorOptions: &pb.WithVectorsSelector_Enable{Enable: false}}, // without vectors
	})
	if err != nil {
		return nil, err
	}

	var mids []primitive.ObjectID
	for idx, p := range res.Result {
		mid := p.GetPayload()["mid"].GetStringValue() // unix timestamp
		mids[idx], err = primitive.ObjectIDFromHex(mid)
		if err != nil {
			return nil, err
		}
	}

	mres, err := ms.mongo.Find(ctx, bson.M{"_id": bson.M{"$in": mids}})
	if err != nil {
		return nil, err
	}
	var memories []*Memory
	err = mres.All(ctx, &memories)
	if err != nil {
		return nil, err
	}
	if len(mids) != len(memories) {
		return nil, fmt.Errorf("some memories not found: \n%v", mids)
	}

	return memories, nil
}

func (ag *Memories) embedding(ctx context.Context, contents []string) ([]openai.Embedding, error) {
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

func (ag *Memories) upsertPoints(ctx context.Context, aid primitive.ObjectID, pids []uuid.UUID, mids []primitive.ObjectID, ems []openai.Embedding) error {
	l := len(ems)
	points := make([]*pb.PointStruct, l)
	for _, em := range ems {
		point := &pb.PointStruct{
			Id:      &pb.PointId{PointIdOptions: &pb.PointId_Uuid{Uuid: pids[em.Index].String()}},
			Payload: map[string]*pb.Value{"mid": {Kind: &pb.Value_StringValue{StringValue: mids[em.Index].Hex()}}},
			Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: em.Embedding}}},
		}
		points[em.Index] = point
	}
	waitUpsert := true
	_, err := ag.qdrant.Upsert(ctx, &pb.UpsertPoints{
		CollectionName: aid.Hex(), // agent's id and also the qdrant collection's name
		Wait:           &waitUpsert,
		Points:         points,
	})

	return err
}

func (ms *Memories) deletePoints(ctx context.Context, aid primitive.ObjectID, pids []uuid.UUID) error {
	ids := make([]*pb.PointId, len(pids))
	for idx, p := range pids {
		ids[idx] = &pb.PointId{PointIdOptions: &pb.PointId_Uuid{Uuid: p.String()}}
	}

	waitDelete := true
	ms.qdrant.Delete(ctx, &pb.DeletePoints{
		CollectionName: aid.Hex(),
		Points:         &pb.PointsSelector{PointsSelectorOneOf: &pb.PointsSelector_Points{Points: &pb.PointsIdsList{Ids: ids}}},
		Wait:           &waitDelete,
	})

	return nil
}
