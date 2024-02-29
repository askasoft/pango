package openai

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/askasoft/pango/log"
)

func testNewAzureOpenAI(t *testing.T) *AzureOpenAI {
	apikey := os.Getenv("AZURE_OPENAI_APIKEY")
	if apikey == "" {
		t.Skip("AZURE_OPENAI_APIKEY not set")
		return nil
	}

	domain := os.Getenv("AZURE_OPENAI_DOMAIN")
	if domain == "" {
		t.Skip("AZURE_OPENAI_DOMAIN not set")
		return nil
	}

	logs := log.NewLog()
	logs.SetLevel(log.LevelDebug)
	aoai := &AzureOpenAI{
		Domain:     domain,
		Apikey:     apikey,
		Apiver:     "2023-05-15",
		Logger:     logs.GetLogger("AZUREOPENAI"),
		MaxRetries: 1,
		RetryAfter: time.Second * 3,
	}

	return aoai
}

func TestAzureOpenAICreateChatCompletion(t *testing.T) {
	aoai := testNewAzureOpenAI(t)
	if aoai == nil {
		return
	}

	aoai.Deployment = os.Getenv("AZURE_OPENAI_CHAT_DEPLOYMENT")

	req := &ChatCompeletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []*ChatMessage{
			{Role: RoleUser, Content: "あなたはだれですか？"},
		},
	}

	res, err := aoai.CreateChatCompletion(req)
	if err != nil {
		t.Fatalf("AzureOpenAI.CreateChatCompletion(): %v", err)
	}

	fmt.Println(res)
}

func TestAzureOpenAICreateTextEmbeddingsAda002(t *testing.T) {
	aoai := testNewAzureOpenAI(t)
	if aoai == nil {
		return
	}

	aoai.Deployment = os.Getenv("AZURE_OPENAI_TEMB_DEPLOYMENT")

	req := &TextEmbeddingsRequest{
		Model: "text-embedding-ada-002",
		Input: []string{"あなたはだれですか？"},
	}

	res, err := aoai.CreateTextEmbeddings(req)
	if err != nil {
		t.Fatalf("AzureOpenAI.CreateTextEmbeddings(): %v", err)
	} else {
		fmt.Println(len(res.Embedding()), res.Usage)
	}
}

func TestAzureOpenAICreateTextEmbeddings3Small(t *testing.T) {
	aoai := testNewAzureOpenAI(t)
	if aoai == nil {
		return
	}

	aoai.Deployment = os.Getenv("AZURE_OPENAI_TEMB_DEPLOYMENT")

	req := &TextEmbeddingsRequest{
		Model: "text-embedding-3-small",
		Input: []string{"あなたはだれですか？"},
	}

	res, err := aoai.CreateTextEmbeddings(req)
	if err != nil {
		t.Fatalf("AzureOpenAI.CreateTextEmbeddings(): %v", err)
	} else {
		fmt.Println(len(res.Embedding()), res.Usage)
	}
}

func TestAzureCreateTextEmbeddings3LargeWithDimensions(t *testing.T) {
	aoai := testNewAzureOpenAI(t)
	if aoai == nil {
		return
	}

	aoai.Deployment = os.Getenv("AZURE_OPENAI_TEMB_DEPLOYMENT")

	req := &TextEmbeddingsRequest{
		Model:      "text-embedding-3-large",
		Input:      []string{"あなたはだれですか？"},
		Dimensions: 1536,
	}

	res, err := aoai.CreateTextEmbeddings(req)
	if err != nil {
		t.Fatalf("AzureOpenAI.CreateTextEmbeddings(): %v", err)
	} else {
		fmt.Println(len(res.Embedding()), res.Usage)
	}
}
