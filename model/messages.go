package model

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Message struct {
	ID             string    `json:"id"`
	UserId         string    `json:"user_id" gorm:"index:idx_user_id_conversation_id_created_at"`
	ConversationId string    `json:"conversation_id" gorm:"index:idx_user_id_conversation_id_created_at"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at" gorm:"index:idx_user_id_conversation_id_created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (m *Message) TableName() string {
	return "messages"
}

func (m *Message) BeforeCreate(db *gorm.DB) error {
	id, err := gonanoid.New()
	if err != nil {
		return err
	}
	m.ID = id
	return nil
}
