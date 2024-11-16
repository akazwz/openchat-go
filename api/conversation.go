package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"openchat/lib"
	"openchat/model"
	"openchat/utils/rcontext"
	"openchat/utils/xhttp"
)

var Conversation = &conversationApi{}

type conversationApi struct{}

func (conversationApi) CreateConversation(w http.ResponseWriter, r *http.Request) {
	userId := rcontext.GetUserId(r.Context())
	conversation := model.Conversation{
		UserId: userId,
	}
	if err := lib.DB.Create(&conversation).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	xhttp.RespJson(w, conversation, http.StatusCreated)
}

func (conversationApi) ListConversations(w http.ResponseWriter, r *http.Request) {
	userId := rcontext.GetUserId(r.Context())
	conversations := make([]*model.Conversation, 0)
	if err := lib.DB.Where("user_id = ?", userId).
		Find(&conversations).
		Order("updated_at DESC").
		Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	xhttp.RespJson(w, conversations, http.StatusOK)
}

func (conversationApi) GetConversation(w http.ResponseWriter, r *http.Request) {
	userId := rcontext.GetUserId(r.Context())
	conversationId := chi.URLParam(r, "id")
	var conversation model.Conversation
	if err := lib.DB.Where("user_id = ? AND id = ?", userId, conversationId).First(&conversation).Error; err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	xhttp.RespJson(w, conversation, http.StatusOK)
}

func (conversationApi) DeleteConversation(w http.ResponseWriter, r *http.Request) {
	userId := rcontext.GetUserId(r.Context())
	conversationId := chi.URLParam(r, "id")
	if err := lib.DB.Where("user_id = ? AND id = ?", userId, conversationId).Delete(model.Conversation{}).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	xhttp.RespJson(w, nil, http.StatusNoContent)
}
