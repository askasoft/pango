package openai

type TextEmbeddingsRequest struct {
	// Input Input text to embed (required)
	Input []string `json:"input,omitempty"`

	// ID of the model to use (required)
	Model string `json:"model,omitempty"`

	// A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse.
	User string `json:"user,omitempty"`
}

func (te *TextEmbeddingsRequest) String() string {
	return toJSONIndent(te)
}

type EmbeddingData struct {
	// The index of the embedding in the list of embeddings.
	Index int `json:"index"`

	// The object type, which is always "embedding".
	Object string `json:"object,omitempty"`

	// The embedding vector, which is a list of floats.
	Embedding []float64 `json:"embedding"`
}

type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

type TextEmbeddingsResponse struct {
	Data   []*EmbeddingData `json:"data"`
	Model  string           `json:"model"`
	Object string           `json:"object"`
	Usage  ChatUsage        `json:"usage"`
}

func (te *TextEmbeddingsResponse) String() string {
	return toJSONIndent(te)
}
