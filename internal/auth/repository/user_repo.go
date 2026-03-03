package repository

import (
	"rate-limiter/internal/auth/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user domain.User) error
	CreateUserPlan(userPlan domain.UserPlan) error
	FindByEmail(email string) (domain.User, error)
	FindByID(id uuid.UUID) (domain.User, error)
	Update(user domain.User) error
	UserExists(email string) (bool, error)
}

type Users struct {
	users     []domain.User
	userPlans []domain.UserPlan
}

func NewUsers() *Users {
	return &Users{
		users: make([]domain.User, 0),
	}
}

func (u *Users) Create(user domain.User) error {
	u.users = append(u.users, user)
	return nil
}

func (u *Users) CreateUserPlan(userPlan domain.UserPlan) error {
	u.userPlans = append(u.userPlans, userPlan)
	return nil
}

func (u *Users) FindByEmail(email string) (domain.User, error) {
	panic("unimplemneted")
}

func (u *Users) UserExists(email string) (bool, error) {
	for _, user := range u.users {
		if user.Email == email {
			return true, nil
		}
	}

	return false, nil
}

func (u *Users) FindByID(id uuid.UUID) (domain.User, error) {
	panic("unimplemneted")
}

func (u *Users) Update(user domain.User) error {
	panic("unimplemneted")
}
