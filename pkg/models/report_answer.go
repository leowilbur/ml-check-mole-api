package models

import (
	"context"

	"github.com/leowilbur/ml-check-mole-api/pkg/types"
)

// ReportAnswer is an answer to a question in a report
type ReportAnswer struct {
	ID         types.ExtendedUUID  `json:"id"`
	ReportID   types.ExtendedUUID  `json:"report_id"`
	QuestionID types.ExtendedUUID  `json:"question_id"`
	Answer     types.ExtendedJSONB `json:"answer"`
}

// CreateReportAnswer creates a report answer
func CreateReportAnswer(ctx context.Context, db Queryer, input *ReportAnswer) error {
	return db.QueryRowEx(
		ctx,
		`INSERT INTO report_answers (
			report_id,
			question_id,
			answer
		) VALUES (
			$1,
			$2,
			$3
		) RETURNING
			id, report_id, question_id, answer`,
		nil,
		&input.ReportID,
		&input.QuestionID,
		&input.Answer,
	).Scan(
		&input.ID,
		&input.ReportID,
		&input.QuestionID,
		&input.Answer,
	)
}

// GetReportAnswer gets a report answer by ID
func GetReportAnswer(
	ctx context.Context, db Queryer, id types.ExtendedUUID,
) (*ReportAnswer, error) {
	item := &ReportAnswer{}
	if err := db.QueryRowEx(
		ctx,
		`SELECT
			id, report_id, question_id, answer
		FROM report_answers WHERE id = $1`,
		nil,
		&id,
	).Scan(
		&item.ID,
		&item.ReportID,
		&item.QuestionID,
		&item.Answer,
	); err != nil {
		return nil, err
	}

	return item, nil
}

// ListReportAnswersByReportID lists all report answers in a report
func ListReportAnswersByReportID(
	ctx context.Context, db Queryer, id types.ExtendedUUID,
) ([]*ReportAnswer, error) {
	rows, err := db.QueryEx(
		ctx,
		`SELECT
			id, report_id, question_id, answer
		FROM report_answers WHERE report_id = $1`,
		nil,
		&id,
	)
	if err != nil {
		return nil, err
	}

	result := []*ReportAnswer{}
	for rows.Next() {
		item := &ReportAnswer{}
		if err := rows.Scan(
			&item.ID,
			&item.ReportID,
			&item.QuestionID,
			&item.Answer,
		); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}

// UpdateReportAnswer updates a single report answer by its ID
func UpdateReportAnswer(ctx context.Context, db Queryer, input *ReportAnswer) error {
	return db.QueryRowEx(
		ctx,
		`UPDATE report_answers SET
			report_id = $2,
			question_id = $3,
			answer = $4
		WHERE id = $1
		RETURNING
			id, report_id, question_id, answer`,
		nil,
		&input.ID,
		&input.ReportID,
		&input.QuestionID,
		&input.Answer,
	).Scan(
		&input.ID,
		&input.ReportID,
		&input.QuestionID,
		&input.Answer,
	)
}

// DeleteReportAnswer deletes a single report answer by its ID
func DeleteReportAnswer(
	ctx context.Context, db Queryer, id types.ExtendedUUID,
) (*ReportAnswer, error) {
	item := &ReportAnswer{}
	if err := db.QueryRowEx(
		ctx,
		`DELETE FROM report_answers WHERE id = $1 RETURNING
			id, report_id, question_id, answer`,
		nil,
		&id,
	).Scan(
		&item.ID,
		&item.ReportID,
		&item.QuestionID,
		&item.Answer,
	); err != nil {
		return nil, err
	}

	return item, nil
}
