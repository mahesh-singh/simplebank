package token

import (
	"testing"
	"time"

	"github.com/mahesh-singh/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	secret := util.RandomString(32)

	maker, err := NewPasetoMaker(secret)
	require.NoError(t, err)
	username := util.RandomOwner()
	duration := time.Hour
	issueAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issueAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestPasetoMarkerExpired(t *testing.T) {
	secret := util.RandomString(32)

	maker, err := NewPasetoMaker(secret)
	require.NoError(t, err)
	username := util.RandomOwner()
	duration := -time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.ErrorContains(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)

}
