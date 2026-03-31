package repository

import (
	"context"

	"kanban/internal/model"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (telegram_id, username, first_name, last_name, photo_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		user.TelegramID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.PhotoURL,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) FindByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, telegram_id, username, first_name, last_name, photo_url, created_at, updated_at
		FROM users
		WHERE telegram_id = $1
	`
	err := r.db.QueryRowContext(ctx, query, telegramID).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.PhotoURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, telegram_id, username, first_name, last_name, photo_url, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.FirstName,
		&user.LastName,
		&user.PhotoURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET username = $1, first_name = $2, last_name = $3, photo_url = $4
		WHERE id = $5
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.PhotoURL,
		user.ID,
	).Scan(&user.UpdatedAt)
}
