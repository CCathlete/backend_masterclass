package token

import (
	"errors"
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

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	}

	// Converting the token string into a jwt object.
	// We need ot provide a pointer to an empty jwt.Claims implementation because the jwt parser will fill it with the payload data.
	jwtToken, err :=
		jwt.ParseWithClaims(signedTokenString, &Payload{}, keyFunc)

	// In case of a invalid token we want to know whether it is because it is invalid or because it has expired.
	if err != nil {
		var validationErr *jwt.ValidationError

		if errors.As(err, &validationErr) {
			// During the parsing, the validationError.Inner field will contain my error coming from my keyFunc call during the parsing.
			if errors.Is(validationErr.Inner, ErrExpiredToken) {
				return nil, ErrExpiredToken
			}
			return nil, ErrInvalidToken
		}
	}

	// jwt.ParseWithClaims returns a jwt.Claims interface, so we need to type assert it to our Payload struct.
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
