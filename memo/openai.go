package memo

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

type ChatMessage struct {
	Role    string `json:"role" bson:"role" toml:"role"`
	Content string `json:"message" bson:"message" toml:"message"`
}

type OpenAI struct {
	client         *openai.Client
	chatModel      string
	emebddingModel openai.EmbeddingModel
}

func NewOpenAI(key string) *OpenAI {
	return &OpenAI{
		client:         openai.NewClient(key),
		chatModel:      openai.GPT3Dot5Turbo,
		emebddingModel: openai.AdaEmbeddingV2,
	}
}

func (oa *OpenAI) Embedding(ctx context.Context, contents []string) (ems []vectors, err error) {
	req := openai.EmbeddingRequest{
		Input: contents,
		Model: openai.AdaEmbeddingV2,
	}

	res, err := oa.client.CreateEmbeddings(ctx, req)
	if err != nil {
		return nil, NewWrapError(500, err, "openai embedding api error occurred")
	}

	ems = make([]vectors, len(contents))
	for _, em := range res.Data {
		ems[em.Index] = em.Embedding // using embedding's own index
	}
	return
}

// Chat is not a streaming chatgpt call
func (oa *OpenAI) Chat(ctx context.Context, messages []ChatMessage) (result ChatMessage, err error) {
	msgs := make([]openai.ChatCompletionMessage, len(messages))
	for i, m := range messages {
		msgs[i] = openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}
	res, err := oa.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    oa.chatModel,
		Messages: msgs,
	})
	if err != nil {
		err = NewWrapError(500, err, "openai chat api error occurred")
		return
	}

	result = ChatMessage{Role: res.Choices[0].Message.Role, Content: res.Choices[0].Message.Content}
	return
}
