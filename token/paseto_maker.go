package token

import (
	tokenUtil "backend-masterclass/token/util"
	"time"

	"golang.org/x/crypto/chacha20poly1305"

	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	*paseto.V2
	secretKey string // Symmetric key.
}

func NewPasetoMaker(secretKey string) (Maker, error) {
	if len(secretKey) != chacha20poly1305.KeySize {
		return nil, tokenUtil.ErrInvalidKeySize
	}

	return &PasetoMaker{
		V2:        &paseto.V2{},
		secretKey: secretKey,
	}, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (signedTokenString string, payload *Payload, err error) {
	payload, err = Newpayload(username, duration)
	if err != nil {
		return
	}

	// ---------Encryted = payload + key + nonce = Signed.----------------
	signedTokenString, err = maker.Encrypt([]byte(maker.secretKey), payload, nil)
	if err != nil {
		return
	}

	return
}

func (maker *PasetoMaker) VerifyToken(signedTokenString string) (payload *Payload, err error) {
	// We need to give a pointer to an empty Payload struct because the paseto parser will fill it with the payload data. In this case the data wiil go directly into our payload unlike jwt where the data was stored in the jwt.Claims interface.
	payload = &Payload{}

	err = maker.Decrypt(signedTokenString, []byte(maker.secretKey), payload, nil)
	if err != nil {
		err = tokenUtil.ErrInvalidToken
		return
	}

	err = payload.Valid()
	if err != nil {
		payload = nil
	}

	return
}
