package model

import (
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID             string    `json:"id"`
	Username       string    `json:"username"`
	HashedPassword string    `json:"-"`
	Nickname       string    `json:"nickname"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Avatar         string    `json:"avatar"`
	Role           Role      `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(db *gorm.DB) error {
	id, err := gonanoid.New()
	if err != nil {
		return err
	}
	u.ID = id
	u.Role = RoleUser
	return nil
}
