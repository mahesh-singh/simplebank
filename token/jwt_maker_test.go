package token

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mahesh-singh/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	secret := util.RandomString(32)
	fmt.Println(secret)
	maker, err := NewJWTMaker(secret)
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

func TestJWTMarkerExpired(t *testing.T) {
	secret := util.RandomString(32)

	maker, err := NewJWTMaker(secret)
	require.NoError(t, err)
	username := util.RandomOwner()
	duration := -time.Minute

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.ErrorContains(t, err, jwt.ErrTokenInvalidClaims.Error())
	require.Nil(t, payload)

}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
		IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
		ID:        payload.ID.String(),
		Audience:  jwt.ClaimStrings{payload.Username},
	})

	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.ErrorContains(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)

}
