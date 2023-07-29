package openai

import (
	"encoding/json"

	"github.com/askasoft/pango/bye"
)

type ChatMessage struct {
	// The role of the messages author. One of system, user, assistant, or function.
	Role string `json:"role,omitempty"`

	// The contents of the message. content is required for all messages, and may be null for assistant messages with function calls.
	Content string `json:"content,omitempty"`

	// The name of the author of this message. name is required if role is function, and it should be the name of the function whose response is in the content. May contain a-z, A-Z, 0-9, and underscores, with a maximum length of 64 characters.
	Name string `json:"name,omitempty"`

	// The name and arguments of a function that should be called, as generated by the model.
	FunctionCall string `json:"function_call,omitempty"`
}

type ChatFunction struct {
	// The name of the function to be called. Must be a-z, A-Z, 0-9, or contain underscores and dashes, with a maximum length of 64.
	Name string `json:"name,omitempty"`

	// A description of what the function does, used by the model to choose when and how to call the function.
	Description string `json:"description,omitempty"`

	// The parameters the functions accepts, described as a JSON Schema object. See the guide for examples, and the JSON Schema reference for documentation about the format.
	// To describe a function that accepts no parameters, provide the value {"type": "object", "properties": {}}.
	Parameters map[string]any `json:"parameters,omitempty"`
}

type ChatCompeletionRequest struct {
	// ID of the model to use. See the model endpoint compatibility table for details on which models work with the Chat API.
	Model string `json:"model,omitempty"`

	// A list of messages comprising the conversation so far.
	Messages []*ChatMessage `json:"messages,omitempty"`

	// A list of functions the model may generate JSON inputs for.
	Functions []*ChatFunction `json:"functions,omitempty"`

	// Controls how the model responds to function calls. "none" means the model does not call a function, and responds to the end-user. "auto" means the model can pick between an end-user or calling a function. Specifying a particular function via {"name":\ "my_function"} forces the model to call that function. "none" is the default when no functions are present. "auto" is the default if functions are present.
	FunctionCall string `json:"function_call,omitempty"`

	// What sampling temperature to use, between 0 and 2. Higher values like 0.8 will make the output more random, while lower values like 0.2 will make it more focused and deterministic.
	// We generally recommend altering this or top_p but not both.
	// Defaults to 1
	Temperature float64 `json:"temperature,omitempty"`

	// An alternative to sampling with temperature, called nucleus sampling, where the model considers the results of the tokens with top_p probability mass. So 0.1 means only the tokens comprising the top 10% probability mass are considered.
	// We generally recommend altering this or temperature but not both.
	// Defaults to 1
	TopP float64 `json:"top_p,omitempty"`

	// How many chat completion choices to generate for each input message.
	// Defaults to 1
	N int `json:"n,omitempty"`

	// If set, partial message deltas will be sent, like in ChatGPT. Tokens will be sent as data-only server-sent events as they become available, with the stream terminated by a data: [DONE] message.
	// Defaults to false
	Stream bool `json:"stream,omitempty"`

	// Up to 4 sequences where the API will stop generating further tokens.
	Stop any `json:"stop,omitempty"`

	// The maximum number of tokens to generate in the chat completion.
	// The total length of input tokens and generated tokens is limited by the model's context length.
	MaxTokens int `json:"max_tokens,omitempty"`

	// Number between -2.0 and 2.0. Positive values penalize new tokens based on whether they appear in the text so far, increasing the model's likelihood to talk about new topics.
	// Defaults to 0
	PresencePenalty float64 `json:"presence_penalty,omitempty"`

	// Number between -2.0 and 2.0. Positive values penalize new tokens based on their existing frequency in the text so far, decreasing the model's likelihood to repeat the same line verbatim.
	// Defaults to 0
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`

	// Modify the likelihood of specified tokens appearing in the completion.
	// Accepts a json object that maps tokens (specified by their token ID in the tokenizer) to an associated bias value from -100 to 100. Mathematically, the bias is added to the logits generated by the model prior to sampling. The exact effect will vary per model, but values between -1 and 1 should decrease or increase likelihood of selection; values like -100 or 100 should result in a ban or exclusive selection of the relevant token.
	// Defaults to null
	LogitBias any `json:"logit_bias,omitempty"`

	// A unique identifier representing your end-user, which can help OpenAI to monitor and detect abuse.
	User string `json:"user,omitempty"`
}

func (cc *ChatCompeletionRequest) String() string {
	bs, err := json.MarshalIndent(cc, "", "  ")
	if err != nil {
		return err.Error()
	}
	return bye.UnsafeString(bs)
}

type ChatChoice struct {
	Index        int         `json:"index,omitempty"`
	Message      ChatMessage `json:"message,omitempty"`
	FinishReason string      `json:"finish_reason,omitempty"`
}

type ChatUsage struct {
	PromtTokens      int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

type ChatCompeletionResponse struct {
	ID      string       `json:"id,omitempty"`
	Object  string       `json:"object,omitempty"`
	Created int64        `json:"created,omitempty"`
	Choices []ChatChoice `json:"choices,omitempty"`
	Usage   *ChatUsage   `json:"usage,omitempty"`
}