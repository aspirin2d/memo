package memo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoFromConfig(t *testing.T) {
	memo := FromConfig("../.config.toml")
	ctx := context.TODO()
	id, err := memo.Agents.Add(ctx, &Agent{Name: "aspirin2d"})
	assert.NoError(t, err)
	err = memo.Agents.Delete(ctx, id)
	assert.NoError(t, err)
}
