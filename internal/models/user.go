package models

import (
	"github.com/satori/uuid"
)

type UserRole string

const (
	RoleClient    UserRole = "client"
	RoleModerator UserRole = "moderator"
)

//easyjson:json
type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Name     string    `json:"name"`
	Role     UserRole  `json:"role"`
}

//easyjson:json
type RegisterRequest struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Name     string   `json:"name"`
	Role     UserRole `json:"role"`
}

//easyjson:json
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//easyjson:json
type LoginResponse struct {
	Token string `json:"token"`
}
