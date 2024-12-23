package token

import (
	u "backend-masterclass/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(u.RandomStr(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := u.RandomUsername()
	duration := time.Minute

	signedTokenString, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, signedTokenString)

	payload, err := maker.VerifyToken(signedTokenString)
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(u.RandomStr(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	token, err := maker.CreateToken(u.RandomUsername(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	// TODO: Add test.
}
