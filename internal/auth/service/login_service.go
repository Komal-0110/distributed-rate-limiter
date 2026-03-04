package service

import (
	"errors"
	"net"
	base "rate-limiter/internal/Base"
	"rate-limiter/internal/auth/domain"
	"rate-limiter/internal/auth/infrastructure"
	"time"

	"github.com/google/uuid"
)

func (s *AuthService) Login(
	email string,
	password string,
	userAgent string,
	ipAddress string,
) (AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return AuthResponse{}, errors.New("invalid credentials")
	}

	if user.Status != domain.Active {
		return AuthResponse{}, errors.New("account is not active")
	}

	cred, err := s.credRepo.FindByUserID(user.ID)
	if err != nil {
		return AuthResponse{}, errors.New("invalid credentials")
	}

	now := time.Now()

	if cred.LockedUntil != nil && cred.LockedUntil.After(now) {
		return AuthResponse{}, errors.New("account is temporarily locked")
	}

	if err := infrastructure.ComparePassword(cred.PasswordHash, password); err != nil {
		if cred.FailedAttempts == 0 {
			failed := 1
			cred.FailedAttempts = failed
		} else {
			cred.FailedAttempts++
		}

		// Lock if too many attempts
		if cred.FailedAttempts >= 5 {
			lockTime := now.Add(15 * time.Minute)
			cred.LockedUntil = &lockTime
		}

		if err := s.credRepo.Update(cred); err != nil {
			return AuthResponse{}, err
		}

		return AuthResponse{}, errors.New("invalid credentials")
	}

	reset := 0
	cred.FailedAttempts = reset
	cred.LockedUntil = nil
	if err := s.credRepo.Update(cred); err != nil {
		return AuthResponse{}, err
	}

	accessToken, err := s.tokenProvider.GenerateAccessToken(
		user.ID,
		user.Email,
	)
	if err != nil {
		return AuthResponse{}, err
	}

	refreshToken, err := s.tokenProvider.GenerateRefreshToken(user.ID)
	if err != nil {
		return AuthResponse{}, err
	}

	session := domain.Session{
		BaseEntity: base.BaseEntity{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		IPAddress:    net.ParseIP(ipAddress),
		IsRevoked:    false,
		ExpiresAt:    now.Add(7 * 24 * time.Hour),
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
