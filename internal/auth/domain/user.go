package domain

import (
	base "rate-limiter/internal/Base"
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	Active    UserStatus = "active"
	Suspended UserStatus = "suspended"
	Deleted   UserStatus = "deleted"
)

type User struct {
	base.BaseEntity

	Email         string
	EmailVerified bool
	Status        UserStatus
}

type Credential struct {
	UserID         uuid.UUID
	PasswordHash   string
	LastChangedAt  time.Time
	FailedAttempts int
	LockedUntil    *time.Time
}

type UserPlanModel string

const (
	Free    UserPlanModel = "free"
	Premium UserPlanModel = "premium"
)

type UserPlan struct {
	UserID    uuid.UUID
	PlanName  UserPlanModel
	ExpiresAt time.Time
}
