package pgstore

import (
	"context"

	"github.com/KozlovNikolai/test-task/internal/model"
)

// CreateUser implements store.IRepository.
func (repo *Repository) CreateUser(context.Context, model.User) (int, error) {
	panic("unimplemented")
}

// DeleteUser implements store.IRepository.
func (repo *Repository) DeleteUser(context.Context, int) error {
	panic("unimplemented")
}

// GetAllUsers implements store.IRepository.
func (repo *Repository) GetAllUsers(ctx context.Context) ([]model.User, error) {
	panic("unimplemented")
}

// GetUserByID implements store.IRepository.
func (repo *Repository) GetUserByID(context.Context, int) (model.User, error) {
	panic("unimplemented")
}

// GetUserByLogin implements store.IRepository.
func (repo *Repository) GetUserByLogin(context.Context, string) (model.User, error) {
	panic("unimplemented")
}
