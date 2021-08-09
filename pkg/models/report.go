package models

import (
	"context"

	"github.com/jackc/pgx/pgtype"
	"github.com/leowilbur/ml-check-mole-api/pkg/types"
)

// Report is a report on the lesion
type Report struct {
	ID                 types.ExtendedUUID        `json:"id"`
	RequestID          types.ExtendedUUID        `json:"request_id"`
	LesionID           types.ExtendedUUID        `json:"lesion_id"`
	Photos             types.ExtendedStringArray `json:"photos"`
	Status             types.ExtendedText        `json:"status"`
	ConsultationResult types.ExtendedText        `json:"consultation_result"`
	CreatedAt          types.ExtendedTimestamp   `json:"created_at"`
	UpdatedAt          types.ExtendedTimestamp   `json:"updated_at"`
}

// CreateReport inserts a report
func CreateReport(ctx context.Context, db Queryer, input *Report) error {
	if input.Status.Status == pgtype.Undefined {
		input.Status.Status = pgtype.Null
	}
	if input.ConsultationResult.Status == pgtype.Undefined {
		input.ConsultationResult.Status = pgtype.Null
	}
	return db.QueryRowEx(
		ctx,
		`INSERT INTO reports (
			request_id,
			lesion_id,
			photos,
			status,
			consultation_result,
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
			id, request_id, lesion_id, photos, status, consultation_result,
			created_at, updated_at`,
		nil,
		&input.RequestID,
		&input.LesionID,
		&input.Photos,
		&input.Status,
		&input.ConsultationResult,
	).Scan(
		&input.ID,
		&input.RequestID,
		&input.LesionID,
		&input.Photos,
		&input.Status,
		&input.ConsultationResult,
		&input.CreatedAt,
		&input.UpdatedAt,
	)
}

// GetReport gets a report
func GetReport(ctx context.Context, db Queryer, id types.ExtendedUUID) (*Report, error) {
	item := &Report{}
	if err := db.QueryRowEx(
		ctx,
		`SELECT
			id, request_id, lesion_id, photos, status, consultation_result,
			created_at, updated_at
		FROM reports WHERE id = $1`,
		nil,
		&id,
	).Scan(
		&item.ID,
		&item.RequestID,
		&item.LesionID,
		&item.Photos,
		&item.Status,
		&item.ConsultationResult,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateReport updates a report
func UpdateReport(ctx context.Context, db Queryer, input *Report) error {
	if input.Status.Status == pgtype.Undefined {
		input.Status.Status = pgtype.Null
	}
	if input.ConsultationResult.Status == pgtype.Undefined {
		input.ConsultationResult.Status = pgtype.Null
	}
	return db.QueryRowEx(
		ctx,
		`UPDATE reports SET
			request_id = $2,
			lesion_id = $3,
			photos = $4,
			status = $5,
			consultation_result = $6,
			updated_at = NOW()
		WHERE id = $1
		RETURNING
			id, request_id, lesion_id, photos, status, consultation_result,
			created_at, updated_at`,
		nil,
		&input.ID,
		&input.RequestID,
		&input.LesionID,
		&input.Photos,
		&input.Status,
		&input.ConsultationResult,
	).Scan(
		&input.ID,
		&input.RequestID,
		&input.LesionID,
		&input.Photos,
		&input.Status,
		&input.ConsultationResult,
		&input.CreatedAt,
		&input.UpdatedAt,
	)
}

// DeleteReport deletes a report
func DeleteReport(ctx context.Context, db Queryer, id types.ExtendedUUID) (*Report, error) {
	item := &Report{}
	if err := db.QueryRowEx(
		ctx,
		`DELETE FROM reports WHERE id = $1 RETURNING
			id, request_id, lesion_id, photos, status, consultation_result,
			created_at, updated_at`,
		nil,
		&id,
	).Scan(
		&item.ID,
		&item.RequestID,
		&item.LesionID,
		&item.Photos,
		&item.Status,
		&item.ConsultationResult,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}
