package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/openai/openai-go"

	"openchat/dto"
	"openchat/lib"
	"openchat/model"
	"openchat/utils/rcontext"
	"openchat/utils/xhttp"
)

var Chat = &chatApi{}

type chatApi struct{}

func (chatApi) ChatCompletion(w http.ResponseWriter, r *http.Request) {
	userId := rcontext.GetUserId(r.Context())
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusBadRequest)
		return
	}
	var reqData dto.ChatReqData
	if err := xhttp.Bind(r, &reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	conversationId := reqData.ConversationId
	messages := reqData.Messages
	params := openai.ChatCompletionNewParams{
		Messages:    openai.F([]openai.ChatCompletionMessageParamUnion{}),
		Model:       openai.F("deepseek-chat"),
		Temperature: openai.F(1.3),
	}
	for _, message := range messages {
		var content any = message.Content
		params.Messages.Value = append(params.Messages.Value, openai.ChatCompletionMessageParam{
			Role:    openai.F(message.Role),
			Content: openai.F(content),
		})
	}
	go func() {
		userMessage := messages[len(messages)-1]
		if userMessage.Role != openai.ChatCompletionMessageParamRoleUser {
			return
		}
		if err := lib.DB.Create(&model.Message{
			UserId:         userId,
			ConversationId: conversationId,
			Role:           string(userMessage.Role),
			Content:        userMessage.Content,
		}); err != nil {
			log.Println(err)
		}
	}()
	stream := lib.DEEPSEEK.Chat.Completions.NewStreaming(context.Background(), params)
	acc := openai.ChatCompletionAccumulator{}
	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)
		if len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			if _, err := w.Write([]byte(content)); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			flusher.Flush()
		}
		if content, ok := acc.JustFinishedContent(); ok {
			log.Println("finished content:", content)
			break
		}
	}
	go func() {
		content := acc.Choices[0].Message.Content
		if err := lib.DB.Create(&model.Message{
			UserId:         userId,
			ConversationId: conversationId,
			Role:           string(openai.ChatCompletionAssistantMessageParamRoleAssistant),
			Content:        content,
		}).Error; err != nil {
			log.Println(err)
		}
		if err := lib.DB.Model(&model.Conversation{}).Where("id = ?", conversationId).Update("updated_at", time.Now()); err != nil {
			log.Println(err)
		}
	}()
	if err := stream.Err(); err != nil {
		log.Println(err)
	}

}

func (chatApi) Summarize(w http.ResponseWriter, r *http.Request) {
	userId := rcontext.GetUserId(r.Context())
	var reqData dto.ChatReqData
	if err := xhttp.Bind(r, &reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	conversationId := reqData.ConversationId
	messages := reqData.Messages
	params := openai.ChatCompletionNewParams{
		Messages:    openai.F([]openai.ChatCompletionMessageParamUnion{}),
		Model:       openai.F("deepseek-chat"),
		Temperature: openai.F(0.0),
	}
	for _, message := range messages {
		var content any = message.Content
		params.Messages.Value = append(params.Messages.Value, openai.ChatCompletionMessageParam{
			Role:    openai.F(message.Role),
			Content: openai.F(content),
		})
	}
	var prompt any = "请你根据上面的对话, 给出一个10个词内的标题,只回答标题即可"
	params.Messages.Value = append(params.Messages.Value, openai.ChatCompletionMessageParam{
		Role:    openai.F(openai.ChatCompletionMessageParamRoleSystem),
		Content: openai.F(prompt),
	})
	completion, err := lib.DEEPSEEK.Chat.Completions.New(context.Background(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	content := completion.Choices[0].Message.Content
	if err = lib.DB.Model(&model.Conversation{}).
		Where("id = ?", conversationId).
		Where("user_id = ?", userId).
		Update("name", content).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	xhttp.RespJson(w, content, http.StatusOK)
}
