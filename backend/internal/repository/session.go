package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"kanban/internal/model"
)

type SessionRepository struct {
	db *DB
}

func NewSessionRepository(db *DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, session *model.Session) error {
	query := `
		INSERT INTO sessions (user_id, token_hash, user_agent, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(ctx, query,
		session.UserID,
		session.TokenHash,
		session.UserAgent,
		session.IPAddress,
		session.ExpiresAt,
	).Scan(&session.ID, &session.CreatedAt)
}

func (r *SessionRepository) FindByToken(ctx context.Context, token string) (*model.Session, error) {
	tokenHash := sha256.Sum256([]byte(token))
	tokenHashHex := hex.EncodeToString(tokenHash[:])
	return r.FindByTokenHash(ctx, tokenHashHex)
}

func (r *SessionRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*model.Session, error) {
	session := &model.Session{}
	query := `
		SELECT id, user_id, token_hash, user_agent, ip_address, created_at, expires_at
		FROM sessions
		WHERE token_hash = $1 AND expires_at > NOW()
	`
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&session.ID,
		&session.UserID,
		&session.TokenHash,
		&session.UserAgent,
		&session.IPAddress,
		&session.CreatedAt,
		&session.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *SessionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID string) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *SessionRepository) DeleteExpired(ctx context.Context) error {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *SessionRepository) FindByUserID(ctx context.Context, userID string) ([]*model.Session, error) {
	query := `
		SELECT id, user_id, token_hash, user_agent, ip_address, created_at, expires_at
		FROM sessions
		WHERE user_id = $1 AND expires_at > NOW()
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]*model.Session, 0)
	for rows.Next() {
		session := &model.Session{}
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.TokenHash,
			&session.UserAgent,
			&session.IPAddress,
			&session.CreatedAt,
			&session.ExpiresAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}
