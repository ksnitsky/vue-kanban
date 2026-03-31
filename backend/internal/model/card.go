package model

import "time"

type Card struct {
	ID        string    `json:"id" db:"id"`
	ColumnID  string    `json:"column_id" db:"column_id"`
	Content   string    `json:"content" db:"content"`
	Position  int       `json:"position" db:"position"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
