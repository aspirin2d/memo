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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Memories is a model which implements MemoryModel interface
// mongo is a mongo collection of memories
// qdrant is a qdrant points of memories
// openai is an openai client
type Memories struct {
	mongo  *mongo.Collection
	qdrant pb.PointsClient
	openai *openai.Client

	SearchLimit int64
	ListLimit   int64
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
	var pids []string = make([]string, l)                         // point ids

	for idx, m := range memories {
		docs[idx] = m
		// check if memory id is nil
		if m.ID != primitive.NilObjectID {
			return nil, NewWrapError(400, fmt.Errorf("memory id should be nil"), "")
		}

		// check if memory aid is nil
		if m.AID != primitive.NilObjectID {
			return nil, NewWrapError(400, fmt.Errorf("agent id should NOT be nil"), "")
		}

		m.ID = primitive.NewObjectID()
		m.AID = aid

		mids[idx] = m.ID
		contents[idx] = m.Content

		// create a reference to the point
		pids[idx] = uuid.New().String()
		m.PID = pids[idx]
	}

	res, err := ms.mongo.InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}

	if len(res.InsertedIDs) != l {
		return nil, NewWrapError(400, fmt.Errorf("some memories not inserted: \n%v\n%v", res.InsertedIDs, mids), "")
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
func (ms *Memories) GetOne(ctx context.Context, aid primitive.ObjectID, mid primitive.ObjectID) (memory *Memory, err error) {
	err = ms.mongo.FindOne(ctx, bson.M{"_id": mid, "aid": aid}).Decode(&memory)
	if err == mongo.ErrNoDocuments {
		return nil, NewWrapError(404, fmt.Errorf("memory not found: %s", mid), "")
	}
	return memory, err
}

// GetMany gets memories by ids
func (ms *Memories) GetMany(ctx context.Context, aid primitive.ObjectID, ids []primitive.ObjectID) (memories []*Memory, err error) {
	cur, err := ms.mongo.Find(ctx, bson.M{"_id": bson.M{"$in": ids}, "aid": aid})
	if err != nil {
		return nil, err
	}
	err = cur.All(ctx, &memories)
	if err != nil {
		return nil, err
	}
	if len(memories) != len(ids) {
		return nil, NewWrapError(400, fmt.Errorf("some memories not found, expected %d but go %d", len(ids), len(memories)), "")
	}
	return
}

// DeleteOne deletes a memory by id
func (ms *Memories) DeleteOne(ctx context.Context, aid primitive.ObjectID, mid primitive.ObjectID) error {
	var mem Memory
	err := ms.mongo.FindOneAndDelete(ctx, bson.M{"_id": mid, "aid": aid}).Decode(&mem)

	if err == mongo.ErrNoDocuments {
		return NewWrapError(404, fmt.Errorf("memory not found: %s", mid), "")
	}

	if err != nil {
		return err
	}

	// also need to delete it from qdrant
	uid, err := uuid.Parse(mem.PID)
	if err != nil {
		return err
	}
	err = ms.deletePoints(ctx, aid, []uuid.UUID{uid})
	return err
}

// DeleteMany deletes memories by ids
func (ms *Memories) DeleteMany(ctx context.Context, aid primitive.ObjectID, ids []primitive.ObjectID) error {
	res, err := ms.mongo.Find(ctx, bson.M{"_id": bson.M{"$in": ids}, "aid": aid})
	if err != nil {
		return err
	}
	var mems []*Memory
	if err = res.All(ctx, &mems); err != nil {
		return err
	}

	// check if all memories found
	if len(mems) != len(ids) {
		return NewWrapError(400, fmt.Errorf("some memories not found"), "")
	}

	// delete memories from mongodb
	_, err = ms.mongo.DeleteMany(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return err
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
	return ms.deletePoints(ctx, aid, pids)
}

// update memory content and its embedding
func (ms *Memories) UpdateOne(ctx context.Context, aid primitive.ObjectID, memory *Memory) error {
	return ms.UpdateMany(ctx, aid, []*Memory{memory})
}

// update memory content and its embedding
func (ms *Memories) UpdateMany(ctx context.Context, aid primitive.ObjectID, memories []*Memory) error {
	upsert := false

	l := len(memories)
	writeModels := make([]mongo.WriteModel, l)
	contents := make([]string, l)
	pids := make([]string, l)
	mids := make([]primitive.ObjectID, l)
	for idx, m := range memories {
		writeModels[idx] = &mongo.UpdateOneModel{Upsert: &upsert, Filter: bson.M{"_id": m.ID, "aid": aid}, Update: bson.M{"$set": bson.M{"content": m.Content}}}
		contents[idx] = m.Content
		pids[idx] = m.PID
		mids[idx] = m.ID
	}

	opts := options.BulkWrite().SetOrdered(false)
	res, err := ms.mongo.BulkWrite(ctx, writeModels, opts)
	if err != nil {
		return err
	}

	// if no memory is modified, then return directly
	if res.ModifiedCount == 0 {
		return NewWrapError(400, fmt.Errorf("memories not modified"), "")
	}

	// generate embedding for new content
	ems, err := ms.embedding(ctx, contents)
	if err != nil {
		return err
	}

	return ms.upsertPoints(ctx, aid, pids, mids, ems)
}

func (ms *Memories) List(ctx context.Context, aid primitive.ObjectID, offset primitive.ObjectID) ([]*Memory, error) {
	filter := bson.M{"aid": aid}
	if offset != primitive.NilObjectID {
		filter["_id"] = bson.M{"$lt": offset} // find memories which are older than offset
	}

	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(ms.ListLimit)
	cursor, err := ms.mongo.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var memories []*Memory
	err = cursor.All(ctx, &memories)
	return memories, err
}

// Search searches memories by query
// id is aeget's id
func (ms *Memories) Search(ctx context.Context, aid primitive.ObjectID, query string) ([]*Memory, []float32, error) {
	ems, err := ms.embedding(ctx, []string{query})
	if err != nil {
		return nil, nil, err
	}
	res, err := ms.qdrant.Search(ctx, &pb.SearchPoints{
		CollectionName: aid.Hex(),
		Vector:         ems[0].Embedding,
		WithPayload:    &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{Enable: true}},  // with payload
		WithVectors:    &pb.WithVectorsSelector{SelectorOptions: &pb.WithVectorsSelector_Enable{Enable: false}}, // without vectors
		Limit:          uint64(ms.SearchLimit),
	})
	if err != nil {
		return nil, nil, err
	}

	// check if any memory found
	if len(res.Result) == 0 {
		return nil, nil, NewWrapError(404, fmt.Errorf("no memories found"), "")
	}

	// get memories from mongodb by ids
	var mids []primitive.ObjectID = make([]primitive.ObjectID, len(res.Result))
	var scores []float32 = make([]float32, len(res.Result))

	for idx, p := range res.Result {
		mid := p.GetPayload()["mid"].GetStringValue() // unix timestamp
		mids[idx], err = primitive.ObjectIDFromHex(mid)
		scores[idx] = p.Score // set scores
		if err != nil {
			return nil, nil, err
		}
	}

	mres, err := ms.mongo.Find(ctx, bson.M{"_id": bson.M{"$in": mids}})
	if err != nil {
		return nil, nil, err
	}

	var memories []*Memory
	err = mres.All(ctx, &memories)
	if err != nil {
		return nil, nil, err
	}

	// check if all memories found
	if len(mids) != len(memories) {
		return nil, nil, NewWrapError(400, fmt.Errorf("some memories not found: \n%v", mids), "")
	}

	// sort by score
	var sorted []*Memory = make([]*Memory, len(memories))
	for _, m := range memories {
		for j, s := range mids {
			if s.Hex() == m.ID.Hex() {
				sorted[j] = m
				// log.Println("sorted: ", j, m.ID.Hex(), m.Content)
				break
			}
		}
	}

	return sorted, scores, nil
}

func (ag *Memories) embedding(ctx context.Context, contents []string) ([]openai.Embedding, error) {
	req := openai.EmbeddingRequest{
		Input: contents,
		Model: openai.AdaEmbeddingV2,
	}

	res, err := ag.openai.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, NewWrapError(500, err, "embedding error")
	}

	return res.Data, nil
}

func (ag *Memories) upsertPoints(ctx context.Context, aid primitive.ObjectID, pids []string, mids []primitive.ObjectID, ems []openai.Embedding) error {
	l := len(ems)
	points := make([]*pb.PointStruct, l)
	for _, em := range ems {
		point := &pb.PointStruct{
			Id:      &pb.PointId{PointIdOptions: &pb.PointId_Uuid{Uuid: pids[em.Index]}},
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

	if err != nil {
		return NewWrapError(500, err, "memory vectors upsert error")
	}

	return nil
}

// deletePoints from qdrant by ids
func (ms *Memories) deletePoints(ctx context.Context, aid primitive.ObjectID, pids []uuid.UUID) error {
	ids := make([]*pb.PointId, len(pids))
	for idx, p := range pids {
		ids[idx] = &pb.PointId{PointIdOptions: &pb.PointId_Uuid{Uuid: p.String()}}
	}

	waitDelete := true
	_, err := ms.qdrant.Delete(ctx, &pb.DeletePoints{
		CollectionName: aid.Hex(),
		Points:         &pb.PointsSelector{PointsSelectorOneOf: &pb.PointsSelector_Points{Points: &pb.PointsIdsList{Ids: ids}}},
		Wait:           &waitDelete,
	})

	if err != nil {
		return NewWrapError(500, err, "memory vectors delete error")
	}
	return nil
}
