package token

import (
	tokenUtil "backend-masterclass/token/util"
	u "backend-masterclass/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	// --------Maker initialisation and token generation--------
	// Creating a new PasetoMaker.
	randomKey := u.RandomStr(32)
	maker, err := NewPasetoMaker(randomKey)
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	// Prepearing data that will be put inside the token's payload.
	username := u.RandomUsername()
	duration := time.Minute

	// Generating a token, returning the payload.
	signedTokenString, beforePayload, err :=
		maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, signedTokenString)
	require.NotNil(t, beforePayload)

	// -----------Token verification, returning the payload---------------
	afterPayload, err := maker.VerifyToken(signedTokenString)
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	require.NoError(t, err)
	require.NotEmpty(t, afterPayload)

	require.NotZero(t, afterPayload.ID)
	require.Equal(t, username, afterPayload.Username)
	require.WithinDuration(t, issuedAt, afterPayload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, afterPayload.ExpiresAt, time.Second)

	// -----------Making sure that the payload is the same----------------
	require.Equal(t, beforePayload, afterPayload)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(u.RandomStr(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	duration := -time.Minute
	token, beforePayload, err := maker.CreateToken(u.RandomUsername(), duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Nil(t, beforePayload)

	afterPayload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, tokenUtil.ErrExpiredToken.Error())
	require.Nil(t, afterPayload)

	require.Equal(t, beforePayload, afterPayload)
}
