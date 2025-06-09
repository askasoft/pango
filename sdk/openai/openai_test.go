package openai

import (
	"context"
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

	req := &ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []*ChatMessage{
			{Role: RoleUser, Content: "あなたはだれですか？"},
		},
	}

	res, err := oai.CreateChatCompletion(context.TODO(), req)
	if err != nil {
		t.Fatalf("OpenAI.CreateChatCompletion(): %v", err)
	}

	fmt.Println("-------------------------------------------")
	fmt.Println(res)
	fmt.Println(res.Usage.String())
}

func TestOpenAIWebSearchTool(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	req := &ChatCompletionRequest{
		Model: "gpt-4o-search-preview",
		Messages: []*ChatMessage{
			{Role: RoleUser, Content: "今年春アニメのおすすめは？"},
		},
		WebSearchOptions: &WebSearchOptions{
			SearchContextSize: "medium",
			UserLocation: &UserLocation{
				Type: "approximate",
				Approximate: &Approximate{
					Country: "JP",
				},
			},
		},
	}

	res, err := oai.CreateChatCompletion(context.TODO(), req)
	if err != nil {
		t.Fatalf("OpenAI.CreateChatCompletion(): %v", err)
	}

	fmt.Println("-------------------------------------------")
	fmt.Println(res)
	fmt.Println(res.Usage.String())
}

func TestOpenAIImageAnalyze(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	req := &ChatCompletionRequest{
		Model: "gpt-4.1",
		Messages: []*ChatMessage{
			{
				Role: RoleUser,
				Content: []*MessageContent{
					{Type: "text", Text: "画像の中に「個人情報が含まれているかどうか」を判定してください。"},
					{Type: "image_url", ImageURL: &ImageURL{
						URL: "https://s3.amazonaws.com/cdn.freshdesk.com/data/helpdesk/attachments/production/50012396079/original/j3UQrTiD9AcapYi98QjFjTKXptsLq4TSBA.png?1720516588",
					}},
				},
			},
		},
	}

	res, err := oai.CreateChatCompletion(context.TODO(), req)
	if err != nil {
		t.Fatalf("OpenAI.CreateChatCompletion(): %v", err)
	}

	fmt.Println("-------------------------------------------")
	fmt.Println(res)
	fmt.Println(res.Usage.String())
}

func TestOpenAICreateTextEmbeddingsAda002(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	req := &TextEmbeddingsRequest{
		Model: "text-embedding-ada-002",
		Input: []string{"あなたはだれですか？"},
	}

	res, err := oai.CreateTextEmbeddings(context.TODO(), req)
	if err != nil {
		t.Fatalf("OpenAI.CreateTextEmbeddings(): %v", err)
	} else {
		fmt.Println(len(res.Embedding()), res.Usage)
	}
}

func TestOpenAICreateTextEmbeddings3Small(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	req := &TextEmbeddingsRequest{
		Model: "text-embedding-3-small",
		Input: []string{"あなたはだれですか？"},
	}

	res, err := oai.CreateTextEmbeddings(context.TODO(), req)
	if err != nil {
		t.Fatalf("OpenAI.CreateTextEmbeddings(): %v", err)
	} else {
		fmt.Println(len(res.Embedding()), res.Usage)
	}
}

func TestOpenAICreateTextEmbeddings3LargeWithDimensions(t *testing.T) {
	oai := testNewOpenAI(t)
	if oai == nil {
		return
	}

	req := &TextEmbeddingsRequest{
		Model:      "text-embedding-3-large",
		Input:      []string{"あなたはだれですか？"},
		Dimensions: 1536,
	}

	res, err := oai.CreateTextEmbeddings(context.TODO(), req)
	if err != nil {
		t.Fatalf("OpenAI.CreateTextEmbeddings(): %v", err)
	} else {
		fmt.Println(len(res.Embedding()), res.Usage)
	}
}
