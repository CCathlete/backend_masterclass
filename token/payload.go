package token

import (
	tokenUtil "backend-masterclass/token/util"
	"time"

	"github.com/google/uuid"
)

// Contains the payload data of the token.
type Payload struct {
	ID        uuid.UUID
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a new token payload with a specific username and duration.
func Newpayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}

// For the payload to implement the jwt.Claims interface we need it to have this method to check if we can create a valid token using this payload.
func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return tokenUtil.ErrExpiredToken
	}
	return nil
}
