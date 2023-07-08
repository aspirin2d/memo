package memo

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type ChatMessage struct {
	Role    string `json:"role" bson:"role" toml:"role"`
	Message string `json:"message" bson:"message" toml:"message"`
}

type OpenAI struct {
	client *openai.Client
}

func NewOpenAI(key string) *OpenAI {
	return &OpenAI{
		client: openai.NewClient(key),
	}
}

func (oa *OpenAI) Embedding(ctx context.Context, contents []string) (ems []vectors, err error) {
	req := openai.EmbeddingRequest{
		Input: contents,
		Model: openai.AdaEmbeddingV2,
	}

	res, err := oa.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, NewWrapError(500, err, "embedding error occurred")
	}

	ems = make([]vectors, len(contents))
	for _, em := range res.Data {
		ems[em.Index] = em.Embedding // using embedding's own index
	}
	return
}

func (oa *OpenAI) Chat(ctx context.Context, messages []ChatMessage) (results []ChatMessage, err error) {
	return
}
