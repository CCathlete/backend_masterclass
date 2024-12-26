package tokenUtil

import "errors"

var (
	ErrInvalidToken   = errors.New("token is invalid")
	ErrExpiredToken   = errors.New("token has expired")
	ErrInvalidKeySize = errors.New("invalid key size")
)
