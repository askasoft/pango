package openai

import (
	"fmt"
	"testing"
)

func TestCreateTextEmbeddings(t *testing.T) {
	openai := testCreateOpenAI()
	if openai == nil {
		return
	}

	te := &TextEmbeddingsRequest{
		Model: "text-embedding-ada-002",
		Input: []string{"You", "We"},
	}

	r, err := openai.CreateTextEmbeddings(te)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println(r.String())
}
