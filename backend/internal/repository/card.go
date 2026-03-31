package repository

import (
	"context"

	"kanban/internal/model"
)

type CardRepository struct {
	db *DB
}

func NewCardRepository(db *DB) *CardRepository {
	return &CardRepository{db: db}
}

func (r *CardRepository) Create(ctx context.Context, card *model.Card) error {
	query := `
		INSERT INTO cards (column_id, content, position)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		card.ColumnID,
		card.Content,
		card.Position,
	).Scan(&card.ID, &card.CreatedAt, &card.UpdatedAt)
}

func (r *CardRepository) FindByID(ctx context.Context, id string) (*model.Card, error) {
	card := &model.Card{}
	query := `
		SELECT id, column_id, content, position, created_at, updated_at
		FROM cards
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&card.ID,
		&card.ColumnID,
		&card.Content,
		&card.Position,
		&card.CreatedAt,
		&card.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return card, nil
}

func (r *CardRepository) FindByColumnID(ctx context.Context, columnID string) ([]*model.Card, error) {
	query := `
		SELECT id, column_id, content, position, created_at, updated_at
		FROM cards
		WHERE column_id = $1
		ORDER BY position ASC
	`
	rows, err := r.db.QueryContext(ctx, query, columnID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cards := make([]*model.Card, 0)
	for rows.Next() {
		card := &model.Card{}
		err := rows.Scan(
			&card.ID,
			&card.ColumnID,
			&card.Content,
			&card.Position,
			&card.CreatedAt,
			&card.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func (r *CardRepository) Update(ctx context.Context, card *model.Card) error {
	query := `
		UPDATE cards
		SET content = $1, position = $2
		WHERE id = $3
		RETURNING updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		card.Content,
		card.Position,
		card.ID,
	).Scan(&card.UpdatedAt)
}

func (r *CardRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM cards WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *CardRepository) GetMaxPosition(ctx context.Context, columnID string) (int, error) {
	query := `
		SELECT COALESCE(MAX(position), 0)
		FROM cards
		WHERE column_id = $1
	`
	var maxPos int
	err := r.db.QueryRowContext(ctx, query, columnID).Scan(&maxPos)
	return maxPos, err
}

func (r *CardRepository) Move(ctx context.Context, cardID string, targetColumnID string, position int) error {
	query := `
		UPDATE cards
		SET column_id = $1, position = $2
		WHERE id = $3
	`
	_, err := r.db.ExecContext(ctx, query, targetColumnID, position, cardID)
	return err
}

func (r *CardRepository) Reorder(ctx context.Context, columnID string, cardIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, id := range cardIDs {
		query := `UPDATE cards SET position = $1 WHERE id = $2 AND column_id = $3`
		_, err := tx.ExecContext(ctx, query, i, id, columnID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
