package model

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Conversation struct {
	ID        string    `json:"id"`
	UserId    string    `json:"user_id" gorm:"index:idx_user_id_updated_at"`
	Name      string    `json:"name"`
	Pinned    bool      `json:"pinned"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"index:idx_user_id_updated_at"`
}

func (c *Conversation) TableName() string {
	return "conversations"
}

func (c *Conversation) BeforeCreate(db *gorm.DB) error {
	id, err := gonanoid.New()
	if err != nil {
		return err
	}
	c.ID = id
	return nil
}
