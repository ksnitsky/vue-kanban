package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"kanban/internal/model"
)

type AuthService struct {
	authTokenRepo AuthTokenRepository
	userRepo      UserRepository
	sessionRepo   SessionRepository
	sessionTTL    time.Duration
}

func NewAuthService(
	authTokenRepo AuthTokenRepository,
	userRepo UserRepository,
	sessionRepo SessionRepository,
	sessionTTL time.Duration,
) *AuthService {
	return &AuthService{
		authTokenRepo: authTokenRepo,
		userRepo:      userRepo,
		sessionRepo:   sessionRepo,
		sessionTTL:    sessionTTL,
	}
}

func (s *AuthService) GenerateToken(ctx context.Context) (*model.AuthToken, error) {
	return s.authTokenRepo.Create(ctx, 5*time.Minute)
}

func (s *AuthService) VerifyToken(ctx context.Context, token string, telegramUser *model.User) (*model.User, error) {
	authToken, err := s.authTokenRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if authToken.Used {
		return nil, ErrTokenUsed
	}

	if authToken.Expired() {
		return nil, ErrTokenExpired
	}

	if err := s.authTokenRepo.MarkUsed(ctx, authToken.ID); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByTelegramID(ctx, telegramUser.TelegramID)
	if err != nil {
		user = &model.User{
			TelegramID: telegramUser.TelegramID,
			Username:   telegramUser.Username,
			FirstName:  telegramUser.FirstName,
			LastName:   telegramUser.LastName,
			PhotoURL:   telegramUser.PhotoURL,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	} else {
		user.Username = telegramUser.Username
		user.FirstName = telegramUser.FirstName
		user.LastName = telegramUser.LastName
		user.PhotoURL = telegramUser.PhotoURL
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *AuthService) CreateSession(ctx context.Context, userID string, userAgent string, ipAddress string) (string, error) {
	token := generateSessionToken()
	tokenHash := sha256.Sum256([]byte(token))

	session := &model.Session{
		UserID:    userID,
		TokenHash: hex.EncodeToString(tokenHash[:]),
		UserAgent: userAgent,
		IPAddress: ipAddress,
		ExpiresAt: time.Now().Add(s.sessionTTL),
	}

	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) ValidateSession(ctx context.Context, token string) (*model.Session, error) {
	session, err := s.sessionRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if session.Expired() {
		return nil, ErrSessionExpired
	}

	return session, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	return s.sessionRepo.Delete(ctx, sessionID)
}

func (s *AuthService) LogoutAll(ctx context.Context, userID string) error {
	return s.sessionRepo.DeleteByUserID(ctx, userID)
}

func (s *AuthService) GetUser(ctx context.Context, userID string) (*model.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

func (s *AuthService) CreateOrGetUser(ctx context.Context, telegramUser *model.User) (*model.User, error) {
	user, err := s.userRepo.FindByTelegramID(ctx, telegramUser.TelegramID)
	if err != nil {
		user = &model.User{
			TelegramID: telegramUser.TelegramID,
			Username:   telegramUser.Username,
			FirstName:  telegramUser.FirstName,
			LastName:   telegramUser.LastName,
			PhotoURL:   telegramUser.PhotoURL,
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	} else {
		user.Username = telegramUser.Username
		user.FirstName = telegramUser.FirstName
		user.LastName = telegramUser.LastName
		user.PhotoURL = telegramUser.PhotoURL
		if err := s.userRepo.Update(ctx, user); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func generateSessionToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

var (
	ErrTokenUsed      = &AuthError{Message: "token already used"}
	ErrTokenExpired   = &AuthError{Message: "token expired"}
	ErrSessionExpired = &AuthError{Message: "session expired"}
)

type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}
