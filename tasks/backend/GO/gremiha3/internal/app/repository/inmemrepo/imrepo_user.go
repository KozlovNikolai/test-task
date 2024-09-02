package inmemrepo

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/KozlovNikolai/test-task/internal/app/repository/models"
)

type UserRepo struct {
	db    *inMemStore
	mutex sync.RWMutex
}

func NewUserRepo(db *inMemStore) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// CreateUser implements services.IUserRepository.
func (repo *UserRepo) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {

	if _, err := repo.GetUserByLogin(ctx, user.Login()); err == nil {
		return domain.User{}, fmt.Errorf("the login %s already exists", user.Login())
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	// мапим домен в модель
	dbUser := domainToUser(user)
	dbUser.ID = repo.db.nextUsersID
	// инкрементируем счетчик записей
	repo.db.nextUsersID++
	// сохраняем
	repo.db.users[dbUser.ID] = dbUser
	log.Printf("modelUser = %v\n", dbUser)
	log.Printf("mapUser = %v\n", repo.db.users[dbUser.ID])

	// мапим модель в домен
	domainUser, err := userToDomain(repo.db.users[dbUser.ID])
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create domain User: %w", err)
	}
	log.Println(domainUser)
	return domainUser, nil
}

// GetUsers implements services.IUserRepository.
func (repo *UserRepo) GetUsers(_ context.Context, limit int, offset int) ([]domain.User, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// извлекаем все ключи из мапы и сортируем их
	keys := make([]int, 0, len(repo.db.users))
	for k := range repo.db.users {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	// выбираем записи с нужными ключами
	var users []models.User
	for i := offset; i < offset+limit && i < len(keys); i++ {
		user := repo.db.users[keys[i]]
		users = append(users, user)
	}

	// мапим массив моделей в массив доменов
	domainUsers := make([]domain.User, len(users))
	for i, user := range users {
		domainUser, err := userToDomain(user)
		if err != nil {
			return nil, fmt.Errorf("failed to create domain User: %w", err)
		}
		domainUsers[i] = domainUser
	}
	return domainUsers, nil
}

// GetUserByID implements services.IUserRepository.
func (repo *UserRepo) GetUserByID(_ context.Context, id int) (domain.User, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	user, exists := repo.db.users[id]
	if !exists {
		return domain.User{}, fmt.Errorf("user with id %d - %w", id, domain.ErrNotFound)
	}
	domainUser, err := userToDomain(user)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create domain User: %w", err)
	}
	return domainUser, nil
}

// GetUsersByOrderID implements services.IUserRepository.
func (repo *UserRepo) GetUserByLogin(_ context.Context, login string) (domain.User, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	var dbUser models.User
	for _, user := range repo.db.users {
		if user.Login == login {
			dbUser = user
			break
		}
	}
	if dbUser.ID == 0 {
		return domain.User{}, fmt.Errorf("user with login %s - %w", login, domain.ErrNotFound)
	}
	domainUser, err := userToDomain(dbUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create domain User: %w", err)
	}
	return domainUser, nil
}

// UpdateUser implements services.IUserRepository.
func (repo *UserRepo) UpdateUser(_ context.Context, user domain.User) (domain.User, error) {
	dbUser := domainToUser(user)
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	// проверяем наличие записи
	_, exists := repo.db.users[dbUser.ID]
	if !exists {
		return domain.User{}, fmt.Errorf("user with id %d - %w", dbUser.ID, domain.ErrNotFound)
	}
	// обновляем запись
	repo.db.users[dbUser.ID] = dbUser
	domainUser, err := userToDomain(dbUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create domain User: %w", err)
	}
	return domainUser, nil
}

// DeleteUser implements services.IUserRepository.
func (repo *UserRepo) DeleteUser(_ context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("%w: id", domain.ErrRequired)
	}
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	_, exists := repo.db.users[id]
	if !exists {
		return fmt.Errorf("user with id %d - %w", id, domain.ErrNotFound)
	}
	delete(repo.db.users, id)
	return nil
}
