package model

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserId    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (rt *RefreshToken) TableName() string {
	return "refresh_tokens"
}

func (rt *RefreshToken) BeforeCreate(db *gorm.DB) error {
	id, err := gonanoid.New()
	if err != nil {
		return err
	}
	rt.ID = id
	return nil
}
