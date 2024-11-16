package model

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Image struct {
	ID         string    `json:"id"`
	UserId     string    `json:"user_id" gorm:"index:idx_user_id_created_at"`
	StorageKey string    `json:"storage_key"`
	Blurhash   string    `json:"blurhash"`
	Prompt     string    `json:"prompt"`
	CreatedAt  time.Time `json:"created_at" gorm:"index:idx_user_id_created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Url string `json:"url" gorm:"-"`
}

func (i *Image) TableName() string {
	return "images"
}

func (i *Image) BeforeCreate(db *gorm.DB) error {
	id, err := gonanoid.New()
	if err != nil {
		return err
	}
	i.ID = id
	return nil
}
