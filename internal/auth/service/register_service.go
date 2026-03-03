package service

import (
	"context"
	"fmt"
	"net/mail"
	base "rate-limiter/internal/Base"
	"rate-limiter/internal/auth/domain"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthResponse struct {
	AccessToken  string
	RefreshToken string
}

func (s *AuthService) Register(ctx context.Context, email string, password string) (AuthResponse, error) {
	if err := s.validateInput(email, password); err != nil {
		return AuthResponse{}, err
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return AuthResponse{}, err
	}

	now := time.Now()
	user := domain.User{
		BaseEntity: base.BaseEntity{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Email:         email,
		EmailVerified: false,
		Status:        domain.Active,
	}

	if err := s.userRepo.Create(user); err != nil {
		return AuthResponse{}, err
	}

	cred := domain.Credential{
		UserID:       user.ID,
		PasswordHash: hashedPassword,
	}

	if err := s.credRepo.Create(cred); err != nil {
		return AuthResponse{}, err
	}

	userPlan := domain.UserPlan{
		UserID:   user.ID,
		PlanName: domain.Free,
	}

	if err := s.userRepo.CreateUserPlan(userPlan); err != nil {
		return AuthResponse{}, err
	}

	accessToken, err := s.tokenProvider.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return AuthResponse{}, err
	}

	refreshToken, err := s.tokenProvider.GenerateRefreshToken(user.ID)
	if err != nil {
		return AuthResponse{}, err
	}

	session := domain.Session{
		BaseEntity: base.BaseEntity{
			ID: uuid.New(),
		},
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return AuthResponse{}, err
	}

	return AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) validateInput(email, password string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email address: %s", email)
	}

	if len(password) < 8 {
		return fmt.Errorf("password len should be greater than 8")
	}

	found, err := s.userRepo.UserExists(email)
	if err != nil {
		return err
	}

	if found {
		return fmt.Errorf("email already exists: %s", email)
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
