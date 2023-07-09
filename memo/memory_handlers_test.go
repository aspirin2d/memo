package memo

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
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
	return []primitive.ObjectID{primitive.NewObjectID()}, mmm.Error
}

func (mmm *mockMemoryModel) GetOne(ctx context.Context, agent primitive.ObjectID, id primitive.ObjectID) (*Memory, error) {
	return &Memory{}, mmm.Error
}
func (mmm *mockMemoryModel) GetMany(ctx context.Context, agent primitive.ObjectID, ids []primitive.ObjectID) ([]*Memory, error) {
	return []*Memory{{Content: "Hello"}, {Content: "World"}}, mmm.Error
}

func (mmm *mockMemoryModel) UpdateOne(ctx context.Context, agent primitive.ObjectID, memory *Memory) error {
	return mmm.Error
}

func (mmm *mockMemoryModel) UpdateMany(ctx context.Context, agent primitive.ObjectID, memories []*Memory) error {
	return mmm.Error
}

func (mmm *mockMemoryModel) DeleteOne(ctx context.Context, agent primitive.ObjectID, id primitive.ObjectID) error {
	return mmm.Error
}
func (mmm *mockMemoryModel) DeleteMany(ctx context.Context, agent primitive.ObjectID, ids []primitive.ObjectID) error {
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

type MemoryHandlersSuite struct {
	suite.Suite
	writer  *httptest.ResponseRecorder
	memo    *Memo
	router  *gin.Engine
	context *gin.Context
}

func (s *MemoryHandlersSuite) SetupSuite() {
	// check memo is implements AgentController
	var _ MemoryController = (*Memo)(nil)

	// create a mock server
	gin.SetMode(gin.ReleaseMode)
	// logger, _ := zap.NewProduction()
	s.memo = &Memo{Memories: &mockMemoryModel{}, Logger: nil}
	s.NotNil(s.memo)
}

func (s *MemoryHandlersSuite) SetupTest() {
	s.writer = httptest.NewRecorder()
	s.context, s.router = gin.CreateTestContext(s.writer)

	s.router.POST("/:aid/add", s.memo.GetAgentId, s.memo.AddMemories)
	s.router.DELETE("/:aid/delete", s.memo.GetAgentId, s.memo.DeleteMemories)
	s.router.PUT("/:aid/update", s.memo.GetAgentId, s.memo.UpdateMemories)
	s.router.GET("/:aid/get", s.memo.GetAgentId, s.memo.GetMemories)
	s.router.GET("/:aid/list", s.memo.GetAgentId, s.memo.ListMemories)
	s.router.GET("/:aid/search", s.memo.GetAgentId, s.memo.SearchMemories)
}
func (s *MemoryHandlersSuite) TearDownTest() {
	s.memo.Memories.(*mockMemoryModel).Error = nil
}

func (s *MemoryHandlersSuite) TestAddMemories() {
	mbody := []map[string]interface{}{
		{"content": "hello, there!"},
	}
	body, _ := json.Marshal(mbody)

	url := "/" + primitive.NewObjectID().Hex() + "/add"
	req := httptest.NewRequest("POST", url, bytes.NewReader(body))
	s.router.ServeHTTP(s.writer, req)
	var m map[string]interface{}
	_ = json.NewDecoder(s.writer.Body).Decode(&m)
	s.NotNil(m["inserted"])
}

func (s *MemoryHandlersSuite) TestDelMemories() {
	mbody := []string{
		new(primitive.ObjectID).Hex(),
		new(primitive.ObjectID).Hex(),
	}
	ids := strings.Join(mbody, ",")
	url := "/" + primitive.NewObjectID().Hex() + "/delete?ids=" + ids
	req := httptest.NewRequest("DELETE", url, nil)
	s.router.ServeHTTP(s.writer, req)

	var m map[string]interface{}
	_ = json.NewDecoder(s.writer.Body).Decode(&m)
	s.NotNil(m["ok"])
}

func (s *MemoryHandlersSuite) TestGetMemories() {
	mbody := []string{
		new(primitive.ObjectID).Hex(),
		new(primitive.ObjectID).Hex(),
	}
	ids := strings.Join(mbody, ",")
	url := "/" + primitive.NewObjectID().Hex() + "/get?ids=" + ids
	req := httptest.NewRequest("GET", url, nil)
	s.router.ServeHTTP(s.writer, req)

	var m []interface{}
	_ = json.NewDecoder(s.writer.Body).Decode(&m)
	s.Equal(len(m), 2)
}

func (s *MemoryHandlersSuite) TestListMemories() {
	url := "/" + primitive.NewObjectID().Hex() + "/list?offset=-1"
	req := httptest.NewRequest("GET", url, nil)
	s.router.ServeHTTP(s.writer, req)

	var m []interface{}
	_ = json.NewDecoder(s.writer.Body).Decode(&m)
	s.Equal(len(m), 5)

	s.writer = httptest.NewRecorder()
	url = "/" + primitive.NewObjectID().Hex() + "/list?offset=" + primitive.NewObjectID().Hex()
	req = httptest.NewRequest("GET", url, nil)
	s.router.ServeHTTP(s.writer, req)

	_ = json.NewDecoder(s.writer.Body).Decode(&m)
	s.Equal(len(m), 5)

	s.writer = httptest.NewRecorder()
	url = "/" + primitive.NewObjectID().Hex() + "/list?offset=abc"
	req = httptest.NewRequest("GET", url, nil)
	s.router.ServeHTTP(s.writer, req)
	s.Equal(s.writer.Code, 400)
}

func (s *MemoryHandlersSuite) TestSearchMemories() {
	url := "/" + primitive.NewObjectID().Hex() + "/search?q=hello"
	req := httptest.NewRequest("GET", url, nil)
	s.router.ServeHTTP(s.writer, req)

	var m map[string]interface{}
	_ = json.NewDecoder(s.writer.Body).Decode(&m)
	s.NotNil(m["memories"])
}

func (s *MemoryHandlersSuite) TestUpdateMemories() {
	mbody := []map[string]interface{}{
		{"id": new(primitive.ObjectID).Hex()},
		{"id": new(primitive.ObjectID).Hex()},
	}
	body, _ := json.Marshal(mbody)
	url := "/" + primitive.NewObjectID().Hex() + "/update"
	req := httptest.NewRequest("PUT", url, bytes.NewReader(body))
	s.router.ServeHTTP(s.writer, req)

	var m map[string]interface{}
	_ = json.NewDecoder(s.writer.Body).Decode(&m)
	s.NotNil(m["ok"])
}

func TestMemoryHandlersSuite(t *testing.T) {
	suite.Run(t, &MemoryHandlersSuite{})
}
