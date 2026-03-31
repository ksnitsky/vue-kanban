package repository

import (
	"context"

	"kanban/internal/model"
)

type BoardRepository struct {
	db *DB
}

func NewBoardRepository(db *DB) *BoardRepository {
	return &BoardRepository{db: db}
}

func (r *BoardRepository) Create(ctx context.Context, board *model.Board) error {
	query := `
		INSERT INTO boards (project_id, name)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		board.ProjectID,
		board.Name,
	).Scan(&board.ID, &board.CreatedAt, &board.UpdatedAt)
}

func (r *BoardRepository) FindByID(ctx context.Context, id string) (*model.Board, error) {
	board := &model.Board{}
	query := `
		SELECT id, project_id, name, created_at, updated_at
		FROM boards
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&board.ID,
		&board.ProjectID,
		&board.Name,
		&board.CreatedAt,
		&board.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return board, nil
}

func (r *BoardRepository) FindByProjectID(ctx context.Context, projectID string) ([]*model.Board, error) {
	query := `
		SELECT id, project_id, name, created_at, updated_at
		FROM boards
		WHERE project_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	boards := make([]*model.Board, 0)
	for rows.Next() {
		board := &model.Board{}
		err := rows.Scan(
			&board.ID,
			&board.ProjectID,
			&board.Name,
			&board.CreatedAt,
			&board.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}
	return boards, nil
}

func (r *BoardRepository) Update(ctx context.Context, board *model.Board) error {
	query := `
		UPDATE boards
		SET name = $1
		WHERE id = $2
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		board.Name,
		board.ID,
	).Scan(&board.UpdatedAt)
}

func (r *BoardRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM boards WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
