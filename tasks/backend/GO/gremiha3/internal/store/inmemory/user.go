package inmemory

import (
	"context"
	"errors"

	"github.com/KozlovNikolai/test-task/internal/model"
)

// Create implements store.IUserRepository.
func (repo *Repository) CreateUser(ctx context.Context, user model.User) (int, error) {

	_, err := repo.GetUserByLogin(ctx, user.Login)
	if err == nil {
		return 0, errors.New("the login already exists")
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	user.ID = repo.nextUsersID
	repo.nextUsersID++
	repo.users[user.ID] = user

	return user.ID, nil
}

// Delete implements store.IUserRepository.
func (i *Repository) DeleteUser(context.Context, int) error {
	panic("unimplemented")
}

// GetAll implements store.IUserRepository.
func (repo *Repository) GetAllUsers(_ context.Context) ([]model.User, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	var users []model.User
	for _, user := range repo.users {
		users = append(users, user)
	}
	return users, nil
}

// GetByID implements store.IUserRepository.
func (repo *Repository) GetUserByID(_ context.Context, userID int) (model.User, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	user, exists := repo.users[userID]
	if !exists {
		return model.User{}, errors.New("user not found")
	}
	return user, nil
}

// GetByLogin implements store.IUserRepository.
func (repo *Repository) GetUserByLogin(_ context.Context, login string) (model.User, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	for _, user := range repo.users {
		if user.Login == login {
			return user, nil
		}
	}
	return model.User{}, errors.New("login not found")
}

// Update implements store.IUserRepository.
func (repo *Repository) UpdateUser(context.Context, int) error {
	panic("unimplemented")
}
