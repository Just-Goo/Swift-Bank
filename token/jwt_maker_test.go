package token

import (
	"testing"
	"time"

	"github.com/Just-Goo/Swift_Bank/helpers"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T)  {
	maker, err := NewJWTMaker(helpers.RandomString(32))
	require.NoError(t, err)

	username := helpers.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.UserName)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}

func TestExpiredJWTToken(t *testing.T)  {
	maker, err := NewJWTMaker(helpers.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(helpers.RandomOwner(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err, ErrExpiredToken)
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T)  {
	payload, err := NewPayload(helpers.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload) // only use for testing
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(helpers.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err, ErrInvalidToken)
	require.Nil(t, payload)
}