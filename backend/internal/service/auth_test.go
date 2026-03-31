package service

import (
	"context"
	"testing"
	"time"

	"kanban/internal/model"
	"kanban/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_GenerateToken(t *testing.T) {
	mockAuthTokenRepo := new(mocks.MockAuthTokenRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockSessionRepo := new(mocks.MockSessionRepository)

	service := NewAuthService(
		mockAuthTokenRepo,
		mockUserRepo,
		mockSessionRepo,
		7*24*time.Hour,
	)

	mockAuthTokenRepo.On("Create", mock.Anything, mock.Anything).Return(&model.AuthToken{
		ID:        "test-id",
		Token:     "test-token",
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}, nil)

	token, err := service.GenerateToken(context.Background())

	assert.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.Token)
	mockAuthTokenRepo.AssertExpectations(t)
}

func TestAuthService_VerifyToken_Success(t *testing.T) {
	mockAuthTokenRepo := new(mocks.MockAuthTokenRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockSessionRepo := new(mocks.MockSessionRepository)

	service := NewAuthService(
		mockAuthTokenRepo,
		mockUserRepo,
		mockSessionRepo,
		7*24*time.Hour,
	)

	authToken := &model.AuthToken{
		ID:        "test-id",
		Token:     "test-token",
		Used:      false,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	telegramUser := &model.User{
		TelegramID: 12345,
		Username:   "testuser",
		FirstName:  "Test",
	}

	mockAuthTokenRepo.On("FindByToken", mock.Anything, "test-token").Return(authToken, nil)
	mockAuthTokenRepo.On("MarkUsed", mock.Anything, "test-id").Return(nil)
	mockUserRepo.On("FindByTelegramID", mock.Anything, int64(12345)).Return(nil, assert.AnError)
	mockUserRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	user, err := service.VerifyToken(context.Background(), "test-token", telegramUser)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, int64(12345), user.TelegramID)
	mockAuthTokenRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_VerifyToken_TokenUsed(t *testing.T) {
	mockAuthTokenRepo := new(mocks.MockAuthTokenRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockSessionRepo := new(mocks.MockSessionRepository)

	service := NewAuthService(
		mockAuthTokenRepo,
		mockUserRepo,
		mockSessionRepo,
		7*24*time.Hour,
	)

	authToken := &model.AuthToken{
		ID:        "test-id",
		Token:     "test-token",
		Used:      true,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	mockAuthTokenRepo.On("FindByToken", mock.Anything, "test-token").Return(authToken, nil)

	user, err := service.VerifyToken(context.Background(), "test-token", &model.User{})

	assert.Error(t, err)
	assert.Equal(t, ErrTokenUsed, err)
	assert.Nil(t, user)
	mockAuthTokenRepo.AssertExpectations(t)
}

func TestAuthService_VerifyToken_TokenExpired(t *testing.T) {
	mockAuthTokenRepo := new(mocks.MockAuthTokenRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockSessionRepo := new(mocks.MockSessionRepository)

	service := NewAuthService(
		mockAuthTokenRepo,
		mockUserRepo,
		mockSessionRepo,
		7*24*time.Hour,
	)

	authToken := &model.AuthToken{
		ID:        "test-id",
		Token:     "test-token",
		Used:      false,
		ExpiresAt: time.Now().Add(-5 * time.Minute),
	}

	mockAuthTokenRepo.On("FindByToken", mock.Anything, "test-token").Return(authToken, nil)

	user, err := service.VerifyToken(context.Background(), "test-token", &model.User{})

	assert.Error(t, err)
	assert.Equal(t, ErrTokenExpired, err)
	assert.Nil(t, user)
	mockAuthTokenRepo.AssertExpectations(t)
}

func TestAuthService_CreateSession(t *testing.T) {
	mockAuthTokenRepo := new(mocks.MockAuthTokenRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockSessionRepo := new(mocks.MockSessionRepository)

	service := NewAuthService(
		mockAuthTokenRepo,
		mockUserRepo,
		mockSessionRepo,
		7*24*time.Hour,
	)

	mockSessionRepo.On("Create", mock.Anything, mock.Anything).Return(nil)

	token, err := service.CreateSession(context.Background(), "user-id", "test-agent", "127.0.0.1")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockSessionRepo.AssertExpectations(t)
}

func TestAuthService_ValidateSession_Success(t *testing.T) {
	mockAuthTokenRepo := new(mocks.MockAuthTokenRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockSessionRepo := new(mocks.MockSessionRepository)

	service := NewAuthService(
		mockAuthTokenRepo,
		mockUserRepo,
		mockSessionRepo,
		7*24*time.Hour,
	)

	session := &model.Session{
		ID:        "session-id",
		UserID:    "user-id",
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	mockSessionRepo.On("FindByToken", mock.Anything, "test-token").Return(session, nil)

	result, err := service.ValidateSession(context.Background(), "test-token")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "user-id", result.UserID)
	mockSessionRepo.AssertExpectations(t)
}

func TestAuthService_ValidateSession_Expired(t *testing.T) {
	mockAuthTokenRepo := new(mocks.MockAuthTokenRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockSessionRepo := new(mocks.MockSessionRepository)

	service := NewAuthService(
		mockAuthTokenRepo,
		mockUserRepo,
		mockSessionRepo,
		7*24*time.Hour,
	)

	session := &model.Session{
		ID:        "session-id",
		UserID:    "user-id",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	mockSessionRepo.On("FindByToken", mock.Anything, "test-token").Return(session, nil)

	result, err := service.ValidateSession(context.Background(), "test-token")

	assert.Error(t, err)
	assert.Equal(t, ErrSessionExpired, err)
	assert.Nil(t, result)
	mockSessionRepo.AssertExpectations(t)
}

func TestAuthService_Logout(t *testing.T) {
	mockAuthTokenRepo := new(mocks.MockAuthTokenRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	mockSessionRepo := new(mocks.MockSessionRepository)

	service := NewAuthService(
		mockAuthTokenRepo,
		mockUserRepo,
		mockSessionRepo,
		7*24*time.Hour,
	)

	mockSessionRepo.On("Delete", mock.Anything, "session-id").Return(nil)

	err := service.Logout(context.Background(), "session-id")

	assert.NoError(t, err)
	mockSessionRepo.AssertExpectations(t)
}
