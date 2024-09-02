package pgrepo

import (
	"context"
	"fmt"

	"github.com/KozlovNikolai/test-task/internal/app/domain"
	"github.com/KozlovNikolai/test-task/internal/app/repository/models"
	"github.com/KozlovNikolai/test-task/internal/pkg/pg"
)

type UserRepo struct {
	db *pg.DB
}

func NewUserRepo(db *pg.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// CreateUser implements services.IUserRepository.
func (u *UserRepo) CreateUser(ctx context.Context, User domain.User) (domain.User, error) {
	dbUser := domainToUser(User)

	var insertedUser models.User

	// Начинаем транзакцию
	tx, err := u.db.WR.Begin(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Вставка данных о пользователе и получение ID
	err = tx.QueryRow(ctx, `
			INSERT INTO users (login,password,role,token)
			VALUES ($1, $2, $3, $4)
			RETURNING id,login,password,role,token`,
		dbUser.Login, dbUser.Password, dbUser.Role, dbUser.Token).
		Scan(
			&insertedUser.ID,
			&insertedUser.Login,
			&insertedUser.Password,
			&insertedUser.Role,
			&insertedUser.Token)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to insert User: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.User{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	domainUser, err := userToDomain(insertedUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create domain User: %w", err)
	}

	return domainUser, nil
}

// DeleteUser implements service.IUserRepository.
func (u *UserRepo) DeleteUser(ctx context.Context, id int) error {
	if id == 0 {
		return fmt.Errorf("%w: id", domain.ErrRequired)
	}
	// Начинаем транзакцию
	tx, err := u.db.WR.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	// Проверяем, что пользователь не связан ни с одним заказом.
	var count int
	err = tx.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM orders
		WHERE user_id = (SELECT id FROM users WHERE id = $1)`, id).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to request the orders users: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("error, there are order-related users.: %w", err)
	}
	// Удаляем пользователя
	_, err = tx.Exec(ctx, `
		DELETE FROM users
		WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete User with id %d: %w", id, err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// GetUsers implements service.IUserRepository.
func (u *UserRepo) GetUsers(ctx context.Context, limit, offset int) ([]domain.User, error) {

	query := `
		SELECT id, login, password, role, token
		FROM users
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	// Запрос
	rows, err := u.db.RO.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	// Заполняем массив пользователей
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Login, &user.Password, &user.Role, &user.Token)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, user)
	}

	// Проверка на ошибки, возникшие при итерации по строкам
	if rows.Err() != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", rows.Err())
	}
	// мапим модель в домен
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

// GetUserByID implements service.IUserRepository.
func (u *UserRepo) GetUserByID(ctx context.Context, id int) (domain.User, error) {
	if id == 0 {
		return domain.User{}, fmt.Errorf("%w: id", domain.ErrRequired)
	}

	// SQL-запрос на получение данных Пользователя по ID
	query := `
		SELECT id, login, password, role, token
		FROM users
		WHERE id = $1
	`
	var user models.User
	// Выполняем запрос и сканируем результат в структуру User
	err := u.db.RO.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Login, &user.Password, &user.Role, &user.Token)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get User by id: %w", err)
	}

	domainUser, err := userToDomain(user)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create domain User: %w", err)
	}

	return domainUser, nil
}

// GetUserByLogin implements service.IUserRepository.
func (u *UserRepo) GetUserByLogin(ctx context.Context, login string) (domain.User, error) {
	if login == "" {
		return domain.User{}, fmt.Errorf("%w: login", domain.ErrRequired)
	}

	// SQL-запрос на получение данных Пользователя по логину
	query := `
		SELECT id, login, password, role, token
		FROM users
		WHERE login = $1
	`
	var user models.User
	// Выполняем запрос и сканируем результат в структуру User
	err := u.db.RO.QueryRow(ctx, query, login).Scan(
		&user.ID, &user.Login, &user.Password, &user.Role, &user.Token)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get User by login: %w", err)
	}

	domainUser, err := userToDomain(user)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create domain User: %w", err)
	}

	return domainUser, nil
}

// UpdateUser implements service.IUserRepository.
func (u *UserRepo) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	dbUser := domainToUser(user)
	//dbUser.UpdatedAt = time.Now()
	// Начинаем транзакцию
	tx, err := u.db.WR.Begin(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	// SQL-запрос на обновление данных Поставщика
	query := `
		UPDATE users
		SET login = $1, password = $2, role = $3, token = $4
		WHERE id = $5
		RETURNING id, login, password, role, token
	`
	var updatedUser models.User

	// Выполняем запрос и сканируем обновленный результат в структуру User
	err = tx.QueryRow(ctx, query,
		dbUser.Login, dbUser.Password, dbUser.Role, dbUser.Token, dbUser.ID).
		Scan(
			&updatedUser.ID,
			&updatedUser.Login,
			&updatedUser.Password,
			&updatedUser.Role,
			&updatedUser.Token)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to update User: %w", err)
	}
	// Фиксация транзакции
	if err := tx.Commit(ctx); err != nil {
		return domain.User{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	domainUser, err := userToDomain(updatedUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create domain User: %w", err)
	}

	return domainUser, nil
}
