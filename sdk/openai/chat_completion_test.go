package openai

import (
	"fmt"
	"os"
	"testing"
)

func testCreateOpenAI() *OpenAI {
	apikey := os.Getenv("OPENAI_API_KEY")
	if apikey == "" {
		return nil
	}

	openai := &OpenAI{
		Domain: "api.openai.com",
		Apikey: apikey,
	}

	return openai
}

func TestCreateChatCompletion(t *testing.T) {
	openai := testCreateOpenAI()
	if openai == nil {
		return
	}

	cc := &ChatCompeletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []*ChatMessage{
			{
				Role:    RoleUser,
				Content: "Who are you?",
			},
		},
	}

	r, err := openai.CreateChatCompletion(cc)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Println(r.String())
}
