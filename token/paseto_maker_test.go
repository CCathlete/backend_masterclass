package token

import (
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

	// Generating a token with an initialised payload.
	signedTokenString, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, signedTokenString)

	// ----------------Token verification------------------------
	payload, err := maker.VerifyToken(signedTokenString)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	// Checking the content of the payload.
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, payload.IssuedAt, time.Now(), time.Second)
	require.WithinDuration(t, payload.ExpiredAt, time.Now().Add(duration), time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(u.RandomStr(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	duration := -time.Minute
	token, err := maker.CreateToken(u.RandomUsername(), duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
