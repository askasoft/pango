package openai

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
)

func testNewOpenAI(t *testing.T) *OpenAI {
	apikey := os.Getenv("OPENAI_APIKEY")
	if apikey == "" {
		t.Skip("OPENAI_APIKEY not set")
		return nil
	}

	logs := log.NewLog()
	logs.SetLevel(log.LevelDebug)
	oai := &OpenAI{
		Domain:     "api.openai.com",
		Apikey:     apikey,
		Logger:     logs.GetLogger("OPENAI"),
		MaxRetries: 1,
		RetryAfter: time.Second * 3,
	}

	return oai
}

func TestOpenAICreateChatCompletion(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	req := &ChatCompeletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []*ChatMessage{
			{Role: RoleUser, Content: "あなたはだれですか？"},
		},
	}

	res, err := oai.CreateChatCompletion(req)
	if err != nil {
		t.Fatalf("OpenAI.CreateChatCompletion(): %v", err)
	}

	fmt.Println(res)
}

func TestOpenAICreateTextEmbeddings(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	req := &TextEmbeddingsRequest{
		Model: "text-embedding-ada-002",
		Input: []string{"あなたはだれですか？"},
	}

	res, err := oai.CreateTextEmbeddings(req)
	if err != nil {
		t.Fatalf("OpenAI.CreateTextEmbeddings(): %v", err)
	}

	fmt.Println(res)
}
