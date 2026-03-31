package repository

import (
	"context"

	"kanban/internal/model"
)

type ColumnRepository struct {
	db *DB
}

func NewColumnRepository(db *DB) *ColumnRepository {
	return &ColumnRepository{db: db}
}

func (r *ColumnRepository) Create(ctx context.Context, column *model.Column) error {
	query := `
		INSERT INTO columns (board_id, title, position)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		column.BoardID,
		column.Title,
		column.Position,
	).Scan(&column.ID, &column.CreatedAt, &column.UpdatedAt)
}

func (r *ColumnRepository) FindByID(ctx context.Context, id string) (*model.Column, error) {
	column := &model.Column{}
	query := `
		SELECT id, board_id, title, position, created_at, updated_at
		FROM columns
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&column.ID,
		&column.BoardID,
		&column.Title,
		&column.Position,
		&column.CreatedAt,
		&column.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return column, nil
}

func (r *ColumnRepository) FindByBoardID(ctx context.Context, boardID string) ([]*model.Column, error) {
	query := `
		SELECT id, board_id, title, position, created_at, updated_at
		FROM columns
		WHERE board_id = $1
		ORDER BY position ASC
	`
	rows, err := r.db.QueryContext(ctx, query, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns := make([]*model.Column, 0)
	for rows.Next() {
		column := &model.Column{}
		err := rows.Scan(
			&column.ID,
			&column.BoardID,
			&column.Title,
			&column.Position,
			&column.CreatedAt,
			&column.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}
	return columns, nil
}

func (r *ColumnRepository) Update(ctx context.Context, column *model.Column) error {
	query := `
		UPDATE columns
		SET title = $1, position = $2
		WHERE id = $3
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		column.Title,
		column.Position,
		column.ID,
	).Scan(&column.UpdatedAt)
}

func (r *ColumnRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM columns WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ColumnRepository) GetMaxPosition(ctx context.Context, boardID string) (int, error) {
	query := `
		SELECT COALESCE(MAX(position), 0)
		FROM columns
		WHERE board_id = $1
	`
	var maxPos int
	err := r.db.QueryRowContext(ctx, query, boardID).Scan(&maxPos)
	return maxPos, err
}

func (r *ColumnRepository) Reorder(ctx context.Context, columnIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, id := range columnIDs {
		query := `UPDATE columns SET position = $1 WHERE id = $2`
		_, err := tx.ExecContext(ctx, query, i, id)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
