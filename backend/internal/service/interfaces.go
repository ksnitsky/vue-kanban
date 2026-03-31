package service

import (
	"context"
	"time"

	"kanban/internal/model"
)

type AuthTokenRepository interface {
	Create(ctx context.Context, ttl time.Duration) (*model.AuthToken, error)
	FindByToken(ctx context.Context, token string) (*model.AuthToken, error)
	MarkUsed(ctx context.Context, id string) error
}

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByTelegramID(ctx context.Context, telegramID int64) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) error
	FindByToken(ctx context.Context, token string) (*model.Session, error)
	Delete(ctx context.Context, id string) error
	DeleteByUserID(ctx context.Context, userID string) error
}
