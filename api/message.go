package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"openchat/lib"
	"openchat/model"
	"openchat/utils/xhttp"
)

var Message = &messageApi{}

type messageApi struct{}

func (messageApi) ListMessages(w http.ResponseWriter, r *http.Request) {
	conversationId := chi.URLParam(r, "id")
	messages := make([]*model.Message, 0)
	if err := lib.DB.Where("conversation_id = ?", conversationId).
		Find(&messages).
		Order("created_at ASC").
		Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	xhttp.RespJson(w, messages, http.StatusOK)
}
