package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (auth *Auth) GetJWTToken(ctx context.Context, creds *Credentials, claims *Claims, currentTime time.Time) (string, error) {
	sessionId, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    "Synctab",
		Subject:   creds.Email,
		ID:        sessionId.String(),
		IssuedAt:  jwt.NewNumericDate(currentTime),
		NotBefore: jwt.NewNumericDate(currentTime),
	}

	expires := currentTime.Add(auth.TokenExpiresDuration)
	claims.ExpiresAt = jwt.NewNumericDate(expires)

	jwtSecret := auth.TokenSecret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
