package models

import (
	"context"

	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/types"
)

// Account is a copy of Cognito's user details
type Account struct {
	ID        types.ExtendedUUID      `json:"id"`
	Name      string                  `json:"name"`
	Email     string                  `json:"email"`
	Phone     string                  `json:"phone"`
	Gender    string                  `json:"gender"`
	BirthDate string                  `json:"birth_date"`
	CreatedAt types.ExtendedTimestamp `json:"created_at"`
	UpdatedAt types.ExtendedTimestamp `json:"updated_at"`
}

// GetAccount fetches a single account
func GetAccount(
	ctx context.Context, db Queryer, id types.ExtendedUUID,
) (*Account, error) {
	account := &Account{}
	if err := db.QueryRowEx(
		ctx,
		`SELECT
			id, name, email, phone, gender,
			birth_date, created_at, updated_at
		FROM accounts WHERE id = $1`,
		nil,
		&id,
	).Scan(
		&account.ID,
		&account.Name,
		&account.Email,
		&account.Phone,
		&account.Gender,
		&account.BirthDate,
		&account.CreatedAt,
		&account.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return account, nil
}

// UpsertAccount upserts the passed account
func UpsertAccount(ctx context.Context, db Queryer, account *Account) error {
	_, err := db.ExecEx(ctx, `INSERT INTO accounts (
		id, name, email, phone, gender, birth_date, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, NOW(), NOW()
	) ON CONFLICT (id) DO UPDATE SET
		name = EXCLUDED.name,
		email = EXCLUDED.email,
		phone = EXCLUDED.phone,
		gender = EXCLUDED.gender,
		birth_date = EXCLUDED.birth_date,
		updated_at = EXCLUDED.created_at
	`, nil,
		&account.ID,
		account.Name,
		account.Email,
		account.Phone,
		account.Gender,
		account.BirthDate,
	)

	return err
}
