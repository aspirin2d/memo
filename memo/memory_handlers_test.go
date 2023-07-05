package memo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

type mockMemoryModel struct {
	Error error
}

func (mmm *mockMemoryModel) AddOne(ctx context.Context, agent primitive.ObjectID, memory *Memory) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), mmm.Error
}
func (mmm *mockMemoryModel) AddMany(ctx context.Context, agent primitive.ObjectID, memories []*Memory) ([]primitive.ObjectID, error) {
	return nil, mmm.Error
}

func (mmm *mockMemoryModel) UpdateOne(ctx context.Context, memory *Memory) error {
	return mmm.Error
}

func (mmm *mockMemoryModel) DeleteOne(ctx context.Context, id primitive.ObjectID) error {
	return mmm.Error
}
func (mmm *mockMemoryModel) DeleteMany(ctx context.Context, ids []primitive.ObjectID) error {
	return mmm.Error
}

func (mmm *mockMemoryModel) List(ctx context.Context, aid primitive.ObjectID, offset primitive.ObjectID) ([]*Memory, error) {
	list := make([]*Memory, 5)
	return list, mmm.Error
}

func (mmm *mockMemoryModel) Search(ctx context.Context, aid primitive.ObjectID, query string) ([]*Memory, []float32, error) {
	list := make([]*Memory, 5)
	scores := make([]float32, 5)
	return list, scores, mmm.Error
}
