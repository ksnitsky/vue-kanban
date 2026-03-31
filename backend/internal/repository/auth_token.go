package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"kanban/internal/model"
)

type AuthTokenRepository struct {
	db *DB
}

func NewAuthTokenRepository(db *DB) *AuthTokenRepository {
	return &AuthTokenRepository{db: db}
}

func (r *AuthTokenRepository) Create(ctx context.Context, ttl time.Duration) (*model.AuthToken, error) {
	token := generateToken(32)
	expiresAt := time.Now().Add(ttl)

	authToken := &model.AuthToken{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	query := `
		INSERT INTO auth_tokens (token, expires_at)
		VALUES ($1, $2)
		RETURNING id, created_at
	`
	err := r.db.QueryRowContext(ctx, query, token, expiresAt).Scan(
		&authToken.ID,
		&authToken.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return authToken, nil
}

func (r *AuthTokenRepository) FindByToken(ctx context.Context, token string) (*model.AuthToken, error) {
	authToken := &model.AuthToken{}
	query := `
		SELECT id, token, used, created_at, expires_at
		FROM auth_tokens
		WHERE token = $1
	`
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&authToken.ID,
		&authToken.Token,
		&authToken.Used,
		&authToken.CreatedAt,
		&authToken.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return authToken, nil
}

func (r *AuthTokenRepository) MarkUsed(ctx context.Context, id string) error {
	query := `UPDATE auth_tokens SET used = TRUE WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *AuthTokenRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM auth_tokens WHERE expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func generateToken(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
