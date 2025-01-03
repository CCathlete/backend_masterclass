package entities

import "time"

type User struct {
	Username          string    `json:"username"`
	PasswordHash      string    `json:"password_hash"`
	Email             string    `json:"email"`
	Role              string    `json:"role"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}
type userRepo interface {
	CreateUser(user User) (User, error)
	GetUser(username string) (User, error)
	ListUsers() ([]User, error)
	UpdateUser(user User) (User, error)
	DeleteUser(username string) error
}
