package memo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmbedding(t *testing.T) {
	memo := FromConfig("../.config.toml")
	oa := NewOpenAI(memo.Config.OpenAIAPIKey)
	ctx := context.TODO()
	ems, err := oa.Embedding(ctx, []string{"hello", "world"})
	assert.NoError(t, err)
	assert.Equal(t, len(ems), 2)
}

func TestChat(t *testing.T) {
	memo := FromConfig("../.config.toml")
	oa := NewOpenAI(memo.Config.OpenAIAPIKey)
	ctx := context.TODO()
	res, err := oa.Chat(ctx, []ChatMessage{{Role: "user", Content: "hello"}})
	assert.NoError(t, err)
	t.Log(res, err)
}
