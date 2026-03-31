package model

import "time"

type AuthToken struct {
	ID        string    `json:"id" db:"id"`
	Token     string    `json:"token" db:"token"`
	Used      bool      `json:"used" db:"used"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
}

func (t *AuthToken) Expired() bool {
	return time.Now().After(t.ExpiresAt)
}
