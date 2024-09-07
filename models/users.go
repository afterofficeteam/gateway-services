package model

import (
	"time"

	"github.com/google/uuid"
)

type UsersLogin struct {
	Id                  uuid.UUID  `json:"id"`
	Email               string     `json:"email"`
	Username            string     `json:"username"`
	Role                string     `json:"role"`
	Address             string     `json:"address"`
	CategoryPreferences []string   `json:"category_preferences"`
	CreatedAt           *time.Time `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at"`
}

type LoginResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
	RefreshToken         string    `json:"refresh_token"`
	RefreshTokenExpiryAt time.Time `json:"refresh_token_expiry_at"`
	*UsersLogin
}
