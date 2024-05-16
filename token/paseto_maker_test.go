package token

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zde37/Swift_Bank/helpers"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPastoMaker(helpers.RandomString(32))
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

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPastoMaker(helpers.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(helpers.RandomOwner(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(helpers.RandomString(10))
	require.Error(t, err, ErrInvalidToken)
	require.Nil(t, payload)
}
