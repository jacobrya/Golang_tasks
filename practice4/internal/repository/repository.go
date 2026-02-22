package repository

import (
	"prac4/internal/repository/_postgres"
	"prac4/internal/repository/_postgres/users"
	"prac4/pkg/modules"
)

type UserRepository interface {
	GetUsers() ([]modules.User, error)
	CreateUser(user *modules.User) (int, error)
	UpdateUser(user *modules.User) error
	GetUserByID(id int) (*modules.User, error)
	DeleteUserByID(id int) (int64, error)
}
type Repositories struct {
	UserRepository
}

func NewRepositories(db *_postgres.Dialect) *Repositories {
	return &Repositories{
		UserRepository: users.NewUserRepository(db),
	}
}
