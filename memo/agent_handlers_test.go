package memo

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type mockAgentsModel struct{}

func (mam *mockAgentsModel) Add(ctx context.Context, agent *Agent) (primitive.ObjectID, error) {
	return primitive.NewObjectID(), nil
}

func (mam *mockAgentsModel) Update(ctx context.Context, agent *Agent) error {
	return nil
}

func (mam *mockAgentsModel) Delete(ctx context.Context, id primitive.ObjectID) error {
	return nil
}

func (mam *mockAgentsModel) Get(ctx context.Context, id primitive.ObjectID) (*Agent, error) {
	return nil, nil
}

func (mam *mockAgentsModel) List(ctx context.Context, offset primitive.ObjectID) ([]*Agent, error) {
	return nil, nil
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

	s.router.POST("/add", s.memo.AddAgent)
}

func (s *AgentHandlersSuite) TestAddAgent() {
	mbody := map[string]interface{}{
		"name": "aspirin2d",
	}
	body, _ := json.Marshal(mbody)

	req := httptest.NewRequest("POST", "/add", bytes.NewReader(body))
	s.router.ServeHTTP(s.writer, req)
	var m map[string]interface{}
	err := json.NewDecoder(s.writer.Body).Decode(&m)
	s.NoError(err)
	_, err = primitive.ObjectIDFromHex(m["inserted"].(string))
	s.NoError(err)
}

func TestAgentHandlersSuite(t *testing.T) {
	suite.Run(t, new(AgentHandlersSuite))
}
