package repository

import (
	"fmt"
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
	for _, cred := range c.credentials {
		if cred.UserID == id {
			return cred, nil
		}
	}

	return domain.Credential{}, fmt.Errorf("not able to find credenetials")
}

func (c *Credentials) Update(cred domain.Credential) error {
	for id, credenetial := range c.credentials {
		if credenetial.UserID == cred.UserID {
			c.credentials[id] = cred
			return nil
		}
	}

	return fmt.Errorf("credentials not found")
}
