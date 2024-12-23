package token

import "time"

type Maker interface {
	CreateToken(username string, duration time.Duration,
	) (signedTokenString string, err error)

	VerifyToken(signedTokenString string) (*Payload, error)
}
