package infrastructure

import (
	"errors"
	"rate-limiter/internal/auth/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type JWTProvider struct {
	secret     []byte
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTProvider(secret string, issuer string) *JWTProvider {
	return &JWTProvider{
		secret:     []byte(secret),
		issuer:     issuer,
		accessTTL:  15 * time.Minute,
		refreshTTL: 7 * 24 * time.Hour,
	}
}

var _ domain.TokenProvider = (*JWTProvider)(nil)

type jwtCustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email,omitempty"`
	Type   string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

func (j *JWTProvider) GenerateAccessToken(
	userID uuid.UUID,
	email string,
) (string, error) {

	claims := jwtCustomClaims{
		UserID: userID.String(),
		Email:  email,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.secret)
}

func (j *JWTProvider) GenerateRefreshToken(
	userID uuid.UUID,
) (string, error) {

	claims := jwtCustomClaims{
		UserID: userID.String(),
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(j.secret)
}

func (j *JWTProvider) parseToken(tokenStr string) (*jwtCustomClaims, error) {

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&jwtCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return j.secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (j *JWTProvider) ValidateAccessToken(
	tokenStr string,
) (*domain.TokenClaims, error) {

	claims, err := j.parseToken(tokenStr)
	if err != nil {
		return nil, err
	}

	if claims.Type != "access" {
		return nil, errors.New("invalid access token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, err
	}

	return &domain.TokenClaims{
		UserID: userID,
		Email:  claims.Email,
	}, nil
}

func (j *JWTProvider) ValidateRefreshToken(
	tokenStr string,
) (*domain.TokenClaims, error) {

	claims, err := j.parseToken(tokenStr)
	if err != nil {
		return nil, err
	}

	if claims.Type != "refresh" {
		return nil, errors.New("invalid refresh token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, err
	}

	return &domain.TokenClaims{
		UserID: userID,
	}, nil
}

func ComparePassword(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(plainPassword),
	)
}
