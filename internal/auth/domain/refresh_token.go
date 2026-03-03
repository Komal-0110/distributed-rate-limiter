package domain

import (
	"net"
	base "rate-limiter/internal/Base"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	base.BaseEntity

	UserID       uuid.UUID
	RefreshToken string
	UserAgent    string
	IPAddress    net.IP
	IsRevoked    bool
	ExpiresAt    time.Time
}

type LoginAttempt struct {
	base.BaseEntity

	Email     string
	IPAddress net.IP
	Success   bool
}

type APIKey struct {
	base.BaseEntity

	UserID   uuid.UUID
	KeyHash  string
	IsActive bool
}

type Role struct {
	base.BaseEntity

	Name string
}

type UserRole struct {
	UserID uuid.UUID
	RoleId uuid.UUID
}
