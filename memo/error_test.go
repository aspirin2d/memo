package memo

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestModelError(t *testing.T) {
	err := NewWrapError(500, nil, "test")
	assert.Equal(t, "test", err.Error())
	assert.Equal(t, err.code, 500)

	w := httptest.NewRecorder()

	gin.SetMode(gin.ReleaseMode)
	ctx := gin.CreateTestContextOnly(w, gin.Default())

	m := &Memo{}
	m.AbortWithError(ctx, err)
	assert.Equal(t, "{\"msg\":\"test\"}", w.Body.String())
}
