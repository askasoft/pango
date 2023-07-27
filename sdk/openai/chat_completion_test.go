package openai

import (
	"fmt"
	"testing"
)

func TestChatCompletionJsonMarshall(t *testing.T) {
	cc := &ChatCompeletionRequest{
		Model: "test",
	}
	fmt.Println(cc.String())
}
