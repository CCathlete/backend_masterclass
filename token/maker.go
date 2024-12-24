package token

import "time"

const (
	// MinSecretKeySize is the minimum size of the secret key.
	MinSecretKeySize = 32
)

type Maker interface {
	CreateToken(username string, duration time.Duration,
	) (signedTokenString string, err error)

	VerifyToken(signedTokenString string) (*Payload, error)
}
