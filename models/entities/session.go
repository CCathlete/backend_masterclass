package entities

import "time"

type Session struct {
	ID           string
	Username     string
	RefreshToken string
	UserAgent    string
	ClientIP     string
	IsBlocked    bool
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

type sessionRepo interface {
	CreateSession(session *Session) error
	GetSession(username string) (*Session, error)
	UpdateSession(session *Session) error
	DeleteSession(token string) error
}
