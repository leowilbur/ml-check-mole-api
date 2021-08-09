package models

import (
	"context"

	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/types"
)

// Question is a question shown to the user when they create a new report
type Question struct {
	ID        types.ExtendedUUID      `db:"id" json:"id"`
	Name      string                  `db:"name" json:"name"`
	Type      string                  `db:"type" json:"type"`
	Answers   types.ExtendedJSONB     `db:"answers" json:"answers"`
	Displayed bool                    `db:"displayed" json:"displayed"`
	Order     int64                   `db:"order" json:"order"`
	CreatedAt types.ExtendedTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt types.ExtendedTimestamp `db:"updated_at" json:"updated_at"`
}

// CreateQuestion inserts a new question
func CreateQuestion(ctx context.Context, db Queryer, input *Question) error {
	return db.QueryRowEx(
		ctx,
		`INSERT INTO questions (
			name,
			type,
			answers,
			displayed,
			"order",
			created_at,
			updated_at
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			NOW(),
			NOW()
		) RETURNING
			id, name, type, answers, displayed,
			"order", created_at, updated_at`,
		nil,
		input.Name,
		input.Type,
		&input.Answers,
		input.Displayed,
		input.Order,
	).Scan(
		&input.ID,
		&input.Name,
		&input.Type,
		&input.Answers,
		&input.Displayed,
		&input.Order,
		&input.CreatedAt,
		&input.UpdatedAt,
	)
}

// GetQuestion gets a single question
func GetQuestion(ctx context.Context, db Queryer, id types.ExtendedUUID) (*Question, error) {
	item := &Question{}
	if err := db.QueryRowEx(
		ctx,
		`SELECT
			id, name, type, answers, displayed,
			"order", created_at, updated_at
		FROM questions WHERE id = $1`,
		nil,
		&id,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Type,
		&item.Answers,
		&item.Displayed,
		&item.Order,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateQuestion updates a single question
func UpdateQuestion(ctx context.Context, db Queryer, input *Question) error {
	return db.QueryRowEx(
		ctx,
		`UPDATE questions SET
			name = $2,
			type = $3,
			answers = $4,
			displayed = $5,
			"order" = $6,
			updated_at = NOW()
		WHERE id = $1
		RETURNING
			id, name, type, answers, displayed, "order",
			created_at, updated_at`,
		nil,
		&input.ID,
		input.Name,
		input.Type,
		&input.Answers,
		input.Displayed,
		input.Order,
	).Scan(
		&input.ID,
		&input.Name,
		&input.Type,
		&input.Answers,
		&input.Displayed,
		&input.Order,
		&input.CreatedAt,
		&input.UpdatedAt,
	)
}

// DeleteQuestion deletes a single qiestion
func DeleteQuestion(ctx context.Context, db Queryer, id types.ExtendedUUID) (*Question, error) {
	item := &Question{}
	if err := db.QueryRowEx(
		ctx,
		`DELETE FROM questions WHERE id = $1 RETURNING
			id, name, type, answers, displayed, "order",
			created_at, updated_at`,
		nil,
		&id,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Type,
		&item.Answers,
		&item.Displayed,
		&item.Order,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}
