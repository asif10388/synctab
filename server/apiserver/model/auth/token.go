package auth

import (
	"context"
	"time"

	"github.com/asif10388/synctab/apiserver/model"
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

func (auth *Auth) ParseJWTToken(ctx context.Context, tokenString string) (*Claims, error) {
	jwtSecret := auth.TokenSecret

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, model.ErrInvalidToken
	}

	return claims, nil
}

func (auth *Auth) ValidateJWTToken(ctx context.Context, tokenString string) (*Claims, error) {
	claims, err := auth.ParseJWTToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
