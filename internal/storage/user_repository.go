package storage

import "github.com/ovk741/TasksStream/internal/domain"

type UserRepository interface {
	Create(user domain.User) error
	GetByEmail(email string) (domain.User, error)
	GetByID(id string) (domain.User, error)
}
