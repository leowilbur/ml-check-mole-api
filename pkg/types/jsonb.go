package types

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/pgtype"
)

// ExtendedJSONB is Postgres JSONB with JSON encoding support
type ExtendedJSONB pgtype.JSONB

// MarshalJSON implements json.Marshaler
func (u *ExtendedJSONB) MarshalJSON() ([]byte, error) {
	switch u.Status {
	case pgtype.Present:
		return u.Bytes, nil
	case pgtype.Null:
		return []byte("null"), nil
	default:
		return nil, fmt.Errorf("Invalid variable status: %+v", u.Status)
	}
}

// UnmarshalJSON implements json.Unmarshaler
func (u *ExtendedJSONB) UnmarshalJSON(input []byte) error {
	*u = ExtendedJSONB{
		Bytes:  input,
		Status: pgtype.Present,
	}

	return nil
}

// Set implements pgx interfaces
func (u *ExtendedJSONB) Set(src interface{}) error {
	return (*pgtype.JSONB)(u).Set(src)
}

// Get implements pgx interfaces
func (u *ExtendedJSONB) Get() interface{} {
	return (*pgtype.JSONB)(u).Get()
}

// AssignTo implements pgx interfaces
func (u *ExtendedJSONB) AssignTo(src interface{}) error {
	return (*pgtype.JSONB)(u).AssignTo(src)
}

// DecodeText implements pgx interfaces
func (u *ExtendedJSONB) DecodeText(ci *pgtype.ConnInfo, src []byte) error {
	return (*pgtype.JSONB)(u).DecodeText(ci, src)
}

// DecodeBinary implements pgx interfaces
func (u *ExtendedJSONB) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	return (*pgtype.JSONB)(u).DecodeBinary(ci, src)
}

// EncodeText implements pgx interfaces
func (u *ExtendedJSONB) EncodeText(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	return (*pgtype.JSONB)(u).EncodeText(ci, buf)
}

// EncodeBinary implements pgx interfaces
func (u *ExtendedJSONB) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	return (*pgtype.JSONB)(u).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (u *ExtendedJSONB) Scan(src interface{}) error {
	return (*pgtype.JSONB)(u).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (u *ExtendedJSONB) Value() (driver.Value, error) {
	return (*pgtype.JSONB)(u).Value()
}
