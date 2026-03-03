package repository

import (
	"rate-limiter/internal/auth/domain"

	"github.com/google/uuid"
)

type CredentialRepository interface {
	Create(cred domain.Credential) error
	FindByUserID(id uuid.UUID) (domain.Credential, error)
	Update(cred domain.Credential) error
}

type Credentials struct {
	credentials []domain.Credential
}

func NewCredentials() *Credentials {
	return &Credentials{
		credentials: make([]domain.Credential, 0),
	}
}

func (c *Credentials) Create(cred domain.Credential) error {
	c.credentials = append(c.credentials, cred)
	return nil
}

func (c *Credentials) FindByUserID(id uuid.UUID) (domain.Credential, error) {
	panic("unimplemneted")
}

func (c *Credentials) Update(cred domain.Credential) error {
	panic("unimplemneted")
}
