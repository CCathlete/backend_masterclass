package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	// MinSecretKeySize is the minimum size of the secret key.
	MinSecretKeySize = 32
)

// JWTMaker is a JSON Web Token maker.
type JWTMaker struct {
	secretKey string
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (signedTokenString string, err error) {
	payload, err := Newpayload(username, duration)
	if err != nil {
		return
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	// We return the string form of the token, including the digital signature.
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// Verifies the validity of the token and returns its payload.
func (maker *JWTMaker) VerifyToken(signedTokenString string,
) (*Payload, error) {

	// The keyFunc is used to validate the signing method of the token and to return the key used to sign the token.
	keyFunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(signedTokenString,
		&Payload{}, keyFunc)
	if err != nil {
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < MinSecretKeySize {
		return nil, ErrInvalidKeySize
	}

	return &JWTMaker{secretKey: secretKey}, nil
}
