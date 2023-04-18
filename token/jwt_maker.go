package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeyLength = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeyLength {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeyLength)
	}
	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	fmt.Printf("payload %v \n", payload)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(payload.ExpiredAt),
		IssuedAt:  jwt.NewNumericDate(payload.IssuedAt),
		ID:        payload.ID.String(),
		Audience:  jwt.ClaimStrings{payload.Username},
	})
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	fmt.Printf("token for parse %v \n", token)

	jwtToken, err := jwt.Parse(token, keyFunc)
	if err != nil {
		fmt.Printf("error while parse %v", err)
		return nil, err
	}

	claims := jwtToken.Claims.(jwt.MapClaims)
	payload := &Payload{
		ID: uuid.MustParse(claims["jti"].(string)),
	}
	if exp, err := claims.GetExpirationTime(); err != nil {
		return nil, fmt.Errorf("error while converting claims exp: %v to payload.ExpiredAt ", claims["exp"])
	} else {
		payload.ExpiredAt = exp.Time
	}
	if iat, err := claims.GetIssuedAt(); err != nil {
		return nil, fmt.Errorf("error while converting claims iat : %v to payload.IssuedAtv", claims["iat"])
	} else {
		payload.IssuedAt = iat.Time
	}
	/*if jit, ok := claims.ge; !ok {
		return nil, fmt.Errorf("error while converting claims jit: %v to payload.ID", claims["jti"])
	} else {
		payload.ID = jit
	}*/
	if aud, err := claims.GetAudience(); err != nil {
		return nil, fmt.Errorf("error while converting claims aud %v to payload.Username", claims["aud"])
	} else {
		payload.Username = aud[0]
	}
	return payload, nil
}
