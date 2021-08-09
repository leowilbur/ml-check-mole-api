package models

import (
	"context"

	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/types"
	"github.com/jackc/pgx/pgtype"
)

// RequestStatus is an ENUM used in the Request model. It describes the current
// state of a request.
type RequestStatus string

var (
	// StatusDraft means that request is just a draft
	StatusDraft RequestStatus = "Draft"
	// StatusOpen means that request has been created
	StatusOpen RequestStatus = "Open"
	// StatusSubmitted means that request has been submitted
	StatusSubmitted RequestStatus = "Submitted"
	// StatusAnswered means that request has been answered
	StatusAnswered RequestStatus = "Answered"
)

// Request is a report review request
type Request struct {
	ID         types.ExtendedUUID      `db:"id" json:"id"`
	AccountID  types.ExtendedUUID      `db:"account" json:"account"`
	Status     *RequestStatus          `db:"status" json:"status"`
	AnswerText types.ExtendedText      `db:"answer_text" json:"answer_text"`
	AnsweredBy types.ExtendedText      `db:"answered_by" json:"answered_by"`
	AnsweredAt types.ExtendedTimestamp `db:"answered_at" json:"answered_at"`
	CreatedAt  types.ExtendedTimestamp `db:"created_at" json:"created_at"`
	UpdatedAt  types.ExtendedTimestamp `db:"updated_at" json:"updated_at"`
}

// CreateRequest creates a request
func CreateRequest(ctx context.Context, db Queryer, input *Request) error {
	if input.AnswerText.Status == pgtype.Undefined {
		input.AnswerText.Status = pgtype.Null
	}
	if input.AnsweredBy.Status == pgtype.Undefined {
		input.AnsweredBy.Status = pgtype.Null
	}
	if input.AnsweredAt.Status == pgtype.Undefined {
		input.AnsweredAt.Status = pgtype.Null
	}

	return db.QueryRowEx(
		ctx,
		`INSERT INTO requests (
			account_id,
			status,
			answer_text,
			answered_by,
			answered_at,
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
			id, account_id, status, answer_text, answered_by, answered_at,
			created_at, updated_at`,
		nil,
		&input.AccountID,
		&input.Status,
		&input.AnswerText,
		&input.AnsweredBy,
		&input.AnsweredAt,
	).Scan(
		&input.ID,
		&input.AccountID,
		&input.Status,
		&input.AnswerText,
		&input.AnsweredBy,
		&input.AnsweredAt,
		&input.CreatedAt,
		&input.UpdatedAt,
	)
}

// GetRequest gets a request
func GetRequest(ctx context.Context, db Queryer, id types.ExtendedUUID) (*Request, error) {
	item := &Request{}
	if err := db.QueryRowEx(
		ctx,
		`SELECT
			id, account_id, status, answer_text, answered_by, answered_at,
			created_at, updated_at
		FROM requests WHERE id = $1`,
		nil,
		&id,
	).Scan(
		&item.ID,
		&item.AccountID,
		&item.Status,
		&item.AnswerText,
		&item.AnsweredBy,
		&item.AnsweredAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}

// UpdateRequest updates a request
func UpdateRequest(ctx context.Context, db Queryer, input *Request) error {
	if input.AnswerText.Status == pgtype.Undefined {
		input.AnswerText.Status = pgtype.Null
	}
	if input.AnsweredBy.Status == pgtype.Undefined {
		input.AnsweredBy.Status = pgtype.Null
	}
	if input.AnsweredAt.Status == pgtype.Undefined {
		input.AnsweredAt.Status = pgtype.Null
	}

	return db.QueryRowEx(
		ctx,
		`UPDATE requests SET
			account_id = $2,
			status = $3,
			answer_text = $4,
			answered_by = $5,
			answered_at = $6,
			updated_at = NOW()
		WHERE id = $1
		RETURNING
			id, account_id, status, answer_text, answered_by, answered_at,
			created_at, updated_at`,
		nil,
		&input.ID,
		&input.AccountID,
		&input.Status,
		&input.AnswerText,
		&input.AnsweredBy,
		&input.AnsweredAt,
	).Scan(
		&input.ID,
		&input.AccountID,
		&input.Status,
		&input.AnswerText,
		&input.AnsweredBy,
		&input.AnsweredAt,
		&input.CreatedAt,
		&input.UpdatedAt,
	)
}

// DeleteRequest deletes a request
func DeleteRequest(ctx context.Context, db Queryer, id types.ExtendedUUID) (*Request, error) {
	item := &Request{}
	if err := db.QueryRowEx(
		ctx,
		`DELETE FROM requests WHERE id = $1 RETURNING
			id, account_id, status, answer_text, answered_by, answered_at,
			created_at, updated_at`,
		nil,
		&id,
	).Scan(
		&item.ID,
		&item.AccountID,
		&item.Status,
		&item.AnswerText,
		&item.AnsweredBy,
		&item.AnsweredAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return item, nil
}
