package dto

import "github.com/openai/openai-go"

type Message struct {
	Role    openai.ChatCompletionMessageParamRole `json:"role"`
	Content string                                `json:"content"`
}

type ChatReqData struct {
	ConversationId string    `json:"conversation_id"`
	Messages       []Message `json:"messages"`
}
