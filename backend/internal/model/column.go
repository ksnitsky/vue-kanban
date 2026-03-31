package model

import "time"

type Column struct {
	ID        string    `json:"id" db:"id"`
	BoardID   string    `json:"board_id" db:"board_id"`
	Title     string    `json:"title" db:"title"`
	Position  int       `json:"position" db:"position"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
