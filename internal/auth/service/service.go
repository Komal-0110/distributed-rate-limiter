package service

import (
	"rate-limiter/internal/auth/domain"
	"rate-limiter/internal/auth/repository"
)

type AuthService struct {
	userRepo      repository.UserRepository
	credRepo      repository.CredentialRepository
	sessionRepo   repository.SessionRepository
	tokenProvider domain.TokenProvider
}

func NewAuthService(
	userRepo repository.UserRepository,
	credRepo repository.CredentialRepository,
	sessionRepo repository.SessionRepository,
	tokenProvider domain.TokenProvider,
) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		credRepo:      credRepo,
		sessionRepo:   sessionRepo,
		tokenProvider: tokenProvider,
	}
}
