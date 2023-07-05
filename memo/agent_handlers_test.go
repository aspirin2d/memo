package memo

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type mockAgentsModel struct {
	Error error // error to return
}

func (mam *mockAgentsModel) Add(ctx context.Context, agent *Agent) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), mam.Error
}

func (mam *mockAgentsModel) Update(ctx context.Context, agent *Agent) error {
	return mam.Error
}

func (mam *mockAgentsModel) Delete(ctx context.Context, id primitive.ObjectID) error {
	return mam.Error
}

func (mam *mockAgentsModel) Get(ctx context.Context, id primitive.ObjectID) (*Agent, error) {
	return &Agent{Name: "aspirin2d"}, mam.Error
}

func (mam *mockAgentsModel) List(ctx context.Context, offset primitive.ObjectID) ([]*Agent, error) {
	list := make([]*Agent, 5)
	return list, nil
}

type AgentHandlersSuite struct {
	suite.Suite
	writer  *httptest.ResponseRecorder
	memo    *Memo
	router  *gin.Engine
	context *gin.Context
}

func (s *AgentHandlersSuite) SetupSuite() {
	// create a mock server
	gin.SetMode(gin.ReleaseMode)
	logger, _ := zap.NewProduction()
	s.memo = &Memo{Agents: &mockAgentsModel{}, Logger: logger.Sugar()}
	s.NotNil(s.memo)
}

func (s *AgentHandlersSuite) SetupTest() {
	s.writer = httptest.NewRecorder()
	s.context, s.router = gin.CreateTestContext(s.writer)

	s.router.GET("/list/:offset", s.memo.ListAgents)
	s.router.GET("/:aid", s.memo.GetAgent)
	s.router.PUT("/add", s.memo.AddAgent)
	s.router.POST("/:aid/update", s.memo.UpdateAgent)
}

func (s *AgentHandlersSuite) TearDownTest() {
	s.memo.Agents.(*mockAgentsModel).Error = nil
}

func (s *AgentHandlersSuite) TestAddAgent() {
	mbody := map[string]interface{}{
		"name": "aspirin2d",
	}
	body, _ := json.Marshal(mbody)

	req := httptest.NewRequest("PUT", "/add", bytes.NewReader(body))
	s.router.ServeHTTP(s.writer, req)
	var m map[string]interface{}
	err := json.NewDecoder(s.writer.Body).Decode(&m)
	s.NoError(err)
	_, err = primitive.ObjectIDFromHex(m["inserted"].(string))
	s.NoError(err)
}

func (s *AgentHandlersSuite) TestAddAgentWithError() {
	s.memo.Agents.(*mockAgentsModel).Error = NewWrapError(400, errors.New("agent id should be nil"), "")

	mbody := map[string]interface{}{
		"id":   primitive.NewObjectID().Hex(),
		"name": "aspirin2d",
	}
	body, _ := json.Marshal(mbody)

	req := httptest.NewRequest("PUT", "/add", bytes.NewReader(body))
	s.router.ServeHTTP(s.writer, req)
	var m map[string]interface{}
	s.Equal(400, s.writer.Code)
	_ = json.NewDecoder(s.writer.Body).Decode(&m)
	s.Equal("agent id should be nil", m["msg"])
}

func (s *AgentHandlersSuite) TestGetAgent() {
	req := httptest.NewRequest("GET", "/123", nil)
	s.router.ServeHTTP(s.writer, req)
	var m map[string]interface{}
	s.Equal(400, s.writer.Code)
	_ = json.NewDecoder(s.writer.Body).Decode(&m)
	s.Equal("invalid agent id", m["msg"])

	s.writer = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/"+primitive.NewObjectID().Hex(), nil)
	s.router.ServeHTTP(s.writer, req)
	s.Equal(200, s.writer.Code)

	var ag Agent
	_ = json.NewDecoder(s.writer.Body).Decode(&ag)
	s.Equal("aspirin2d", ag.Name)
}

func (s *AgentHandlersSuite) TestUpdateAgent() {
	mbody := map[string]interface{}{
		"name": "aspirin2d",
	}
	body, _ := json.Marshal(mbody)
	req := httptest.NewRequest("POST", "/"+primitive.NewObjectID().Hex()+"/update", bytes.NewReader(body))
	s.router.ServeHTTP(s.writer, req)
	var m map[string]interface{}
	s.Equal(200, s.writer.Code)
	_ = json.NewDecoder(s.writer.Body).Decode(&m)
	s.Equal(true, m["ok"])
}

func (s *AgentHandlersSuite) TestListAgents() {
	req := httptest.NewRequest("GET", "/list/nil", nil)
	s.router.ServeHTTP(s.writer, req)
	var m []interface{}
	s.Equal(200, s.writer.Code)
	_ = json.NewDecoder(s.writer.Body).Decode(&m)
	s.Equal(5, len(m))

	s.writer = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/list/-1", nil)
	s.router.ServeHTTP(s.writer, req)
	s.Equal(200, s.writer.Code)
	_ = json.NewDecoder(s.writer.Body).Decode(&m)
	s.Equal(5, len(m))
}

func TestAgentHandlersSuite(t *testing.T) {
	suite.Run(t, new(AgentHandlersSuite))
}
