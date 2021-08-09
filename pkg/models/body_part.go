package models

import (
	"context"

	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/types"
	"github.com/jackc/pgx/pgtype"
)

// BodyPart contains information about a body part shown in the UI
type BodyPart struct {
	ID        types.ExtendedUUID `db:"id" json:"id"`
	Name      string             `db:"name" json:"name"`
	Displayed bool               `db:"displayed" json:"displayed"`
	Image     string             `db:"image" json:"image"`
	Order     int                `db:"order" json:"order"`
	Parent    types.ExtendedUUID `db:"parent" json:"parent"`
}

// CreateBodyPart inserts a single body part into the store
func CreateBodyPart(ctx context.Context, db Queryer, input *BodyPart) error {
	if input.Parent.Status == pgtype.Undefined {
		input.Parent.Status = pgtype.Null
	}

	return db.QueryRowEx(
		ctx,
		`INSERT INTO body_parts (
			name,
			displayed,
			image,
			"order",
			parent
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		) RETURNING
			id, name, displayed, image, "order", parent`,
		nil,
		input.Name,
		input.Displayed,
		input.Image,
		input.Order,
		&input.Parent,
	).Scan(
		&input.ID,
		&input.Name,
		&input.Displayed,
		&input.Image,
		&input.Order,
		&input.Parent,
	)
}

// GetBodyPart gets a single body part from the database by its ID
func GetBodyPart(ctx context.Context, db Queryer, id types.ExtendedUUID) (*BodyPart, error) {
	item := &BodyPart{}
	if err := db.QueryRowEx(
		ctx,
		`SELECT
			id, name, displayed, image, "order", parent
		FROM body_parts WHERE id = $1`,
		nil,
		&id,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Displayed,
		&item.Image,
		&item.Order,
		&item.Parent,
	); err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateBodyPart updates the passed body part
func UpdateBodyPart(ctx context.Context, db Queryer, input *BodyPart) error {
	if input.Parent.Status == pgtype.Undefined {
		input.Parent.Status = pgtype.Null
	}

	return db.QueryRowEx(
		ctx,
		`UPDATE body_parts SET
			name = $2,
			displayed = $3,
			image = $4,
			"order" = $5,
			parent = $6
		WHERE id = $1
		RETURNING
			id, name, displayed, image, "order", parent`,
		nil,
		&input.ID,
		input.Name,
		input.Displayed,
		input.Image,
		input.Order,
		&input.Parent,
	).Scan(
		&input.ID,
		&input.Name,
		&input.Displayed,
		&input.Image,
		&input.Order,
		&input.Parent,
	)
}

// DeleteBodyPart deletes a body part by its ID and returns it
func DeleteBodyPart(ctx context.Context, db Queryer, id types.ExtendedUUID) (*BodyPart, error) {
	item := &BodyPart{}
	if err := db.QueryRowEx(
		ctx,
		`DELETE FROM body_parts WHERE id = $1 RETURNING id, name, displayed, image, "order", parent`,
		nil,
		&id,
	).Scan(
		&item.ID,
		&item.Name,
		&item.Displayed,
		&item.Image,
		&item.Order,
		&item.Parent,
	); err != nil {
		return nil, err
	}

	return item, nil
}
