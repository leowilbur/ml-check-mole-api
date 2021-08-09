package models

import (
	"context"

	"github.com/leowilbur/ml-check-mole-api/pkg/types"
)

// Lesion describes user's lesion
type Lesion struct {
	ID               types.ExtendedUUID      `json:"id"`
	AccountID        types.ExtendedUUID      `json:"account_id"`
	Name             string                  `json:"name"`
	BodyPartID       types.ExtendedUUID      `json:"body_part_id"`
	BodyPartLocation string                  `json:"body_part_location"`
	CreatedAt        types.ExtendedTimestamp `json:"created_at"`
	UpdatedAt        types.ExtendedTimestamp `json:"updated_at"`
}

// CreateLesion inserts a new lession
func CreateLesion(ctx context.Context, db Queryer, input *Lesion) error {
	return db.QueryRowEx(
		ctx,
		`INSERT INTO lesions (
			account_id,
			name,
			body_part_id,
			body_part_location,
			created_at,
			updated_at
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			NOW(),
			NOW()
		) RETURNING
			id, account_id, name, body_part_id, body_part_location,
			created_at, updated_at`,
		nil,
		&input.AccountID,
		input.Name,
		&input.BodyPartID,
		input.BodyPartLocation,
	).Scan(
		&input.ID,
		&input.AccountID,
		&input.Name,
		&input.BodyPartID,
		&input.BodyPartLocation,
		&input.CreatedAt,
		&input.UpdatedAt,
	)
}

// GetLesion fetches a single lesion
func GetLesion(ctx context.Context, db Queryer, id types.ExtendedUUID) (*Lesion, error) {
	item := &Lesion{}
	if err := db.QueryRowEx(
		ctx,
		`SELECT
			id, account_id, name, body_part_id, body_part_location,
			created_at, updated_at
		FROM lesions WHERE id = $1`,
		nil,
		&id,
	).Scan(
		&item.ID,
		&item.AccountID,
		&item.Name,
		&item.BodyPartID,
		&item.BodyPartLocation,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateLesion updates the specified lesion
func UpdateLesion(ctx context.Context, db Queryer, input *Lesion) error {
	return db.QueryRowEx(
		ctx,
		`UPDATE lesions SET
			account_id = $2,
			name = $3,
			body_part_id = $4,
			body_part_location = $5,
			updated_at = NOW()
		WHERE id = $1
		RETURNING
			id, account_id, name, body_part_id, body_part_location,
			created_at, updated_at`,
		nil,
		&input.ID,
		&input.AccountID,
		input.Name,
		&input.BodyPartID,
		input.BodyPartLocation,
	).Scan(
		&input.ID,
		&input.AccountID,
		&input.Name,
		&input.BodyPartID,
		&input.BodyPartLocation,
		&input.CreatedAt,
		&input.UpdatedAt,
	)
}

// DeleteLesion deletes a single lesion.
func DeleteLesion(ctx context.Context, db Queryer, id types.ExtendedUUID) (*Lesion, error) {
	item := &Lesion{}
	if err := db.QueryRowEx(
		ctx,
		`DELETE FROM lesions WHERE id = $1 RETURNING
			id, account_id, name, body_part_id, body_part_location,
			created_at, updated_at`,
		nil,
		&id,
	).Scan(
		&item.ID,
		&item.AccountID,
		&item.Name,
		&item.BodyPartID,
		&item.BodyPartLocation,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}
