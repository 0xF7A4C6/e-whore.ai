package vercel

import "github.com/Implex-ltd/cleanhttp/cleanhttp"

type Client struct {
	Http        *cleanhttp.CleanHttp
	ChatUUID    string
	Model       string
	Temperature float64
	MaxTokens   int
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SendPromptPayload struct {
	Messages         []Message `json:"messages"`
	PlaygroundID     string    `json:"playgroundId"`
	ChatIndex        int       `json:"chatIndex"`
	Model            string    `json:"model"`
	Temperature      float64   `json:"temperature"`
	MaxTokens        int       `json:"maxTokens"`
	TopK             int       `json:"topK"`
	TopP             int       `json:"topP"`
	FrequencyPenalty int       `json:"frequencyPenalty"`
	PresencePenalty  int       `json:"presencePenalty"`
	StopSequences    []string  `json:"stopSequences"`
}

type B64Payload struct {
	T string  `json:"t"`
	C string  `json:"c"`
	A float64 `json:"a"`
}