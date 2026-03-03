package domain

import "github.com/google/uuid"

type TokenProvider interface {
	GenerateAccessToken(userID uuid.UUID, email string) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)

	ValidateAccessToken(token string) (*TokenClaims, error)
	ValidateRefreshToken(token string) (*TokenClaims, error)
}

type TokenClaims struct {
	UserID uuid.UUID
	Email  string
}
