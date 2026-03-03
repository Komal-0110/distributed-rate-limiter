package repository

import (
	"rate-limiter/internal/auth/domain"

	"github.com/google/uuid"
)

type SessionRepository interface {
	Create(session domain.Session) error
	FindByRefreshToken(token string) (domain.Session, error)
	Revoke(sessionID uuid.UUID) error
	RevokeAllByUser(userID uuid.UUID) error
}

type Sessions struct {
	sessions []domain.Session
}

func NewSessions() *Sessions {
	return &Sessions{
		sessions: make([]domain.Session, 0),
	}
}

func (s *Sessions) Create(session domain.Session) error {
	s.sessions = append(s.sessions, session)
	return nil
}

func (s *Sessions) FindByRefreshToken(token string) (domain.Session, error) {
	panic("unimplemneted")
}

func (s *Sessions) Revoke(sessionID uuid.UUID) error {
	panic("unimplemneted")
}

func (s *Sessions) RevokeAllByUser(userID uuid.UUID) error {
	panic("unimplemneted")
}
