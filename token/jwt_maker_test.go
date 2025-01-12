package token

import (
	tokenUtil "backend-masterclass/token/util"
	u "backend-masterclass/util"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(u.RandomStr(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := u.RandomUsername()
	duration := time.Minute

	// -----------Token generation, returning the payload-----------------
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

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(u.RandomStr(32))
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

// Testing for use of alg "none".
// We want to make sure that we get an ErrInvalidToken error when we get a token that was signed with the none alg.
func TestInvalidJWTTokenAlgNone(t *testing.T) {

	// ---------This part imitates maker.CreateToken()------------
	payload, err := Newpayload(u.RandomUsername(), time.Minute)
	require.NoError(t, err)

	// Creating a token object with our payload, putting in the none Alg using the signing method none variable.
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	// We need to use a special key to be able to sign the token with alg "none".
	junkKey := jwt.UnsafeAllowNoneSignatureType
	signedTokenString, err := jwtToken.SignedString(junkKey)
	require.NoError(t, err)

	// ---------This part imitates maker.VerifyToken()------------

	// Creating a new JWTMaker with a our "real" key.
	myRealKey := u.RandomStr(32)
	maker, err := NewJWTMaker(myRealKey)
	require.NoError(t, err)

	// Now for the big part, we want to test our verification method, showing that it won't verify a token signed with a different key even if the signing algorithm was none.
	payload, err = maker.VerifyToken(signedTokenString)
	require.Error(t, err)
	require.EqualError(t, err, tokenUtil.ErrInvalidToken.Error())
	require.Nil(t, payload)
}
