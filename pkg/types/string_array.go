package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/pgtype"
)

// ExtendedStringArray is an array of varchars with JSON encoding support
type ExtendedStringArray pgtype.VarcharArray

// MarshalJSON implements json.Marshaler
func (u *ExtendedStringArray) MarshalJSON() ([]byte, error) {
	switch u.Status {
	case pgtype.Present:
		target := []string{}
		if err := u.AssignTo(&target); err != nil {
			return nil, err
		}

		return json.Marshal(target)
	case pgtype.Null:
		return []byte("null"), nil
	default:
		return nil, fmt.Errorf("Invalid variable status: %+v", u.Status)
	}
}

// UnmarshalJSON implements json.Unmarshaler
func (u *ExtendedStringArray) UnmarshalJSON(input []byte) error {
	if bytes.Equal(input, []byte("null")) {
		u = &ExtendedStringArray{
			Status: pgtype.Null,
		}
		return nil
	}

	*u = ExtendedStringArray{}

	result := []string{}
	if err := json.Unmarshal(input, &result); err != nil {
		return err
	}
	if err := u.Set(result); err != nil {
		return err
	}

	return nil
}

// Set implements pgx interfaces
func (u *ExtendedStringArray) Set(src interface{}) error {
	return (*pgtype.VarcharArray)(u).Set(src)
}

// Get implements pgx interfaces
func (u *ExtendedStringArray) Get() interface{} {
	return (*pgtype.VarcharArray)(u).Get()
}

// AssignTo implements pgx interfaces
func (u *ExtendedStringArray) AssignTo(src interface{}) error {
	return (*pgtype.VarcharArray)(u).AssignTo(src)
}

// DecodeText implements pgx interfaces
func (u *ExtendedStringArray) DecodeText(ci *pgtype.ConnInfo, src []byte) error {
	return (*pgtype.VarcharArray)(u).DecodeText(ci, src)
}

// DecodeBinary implements pgx interfaces
func (u *ExtendedStringArray) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	return (*pgtype.VarcharArray)(u).DecodeBinary(ci, src)
}

// EncodeText implements pgx interfaces
func (u *ExtendedStringArray) EncodeText(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	return (*pgtype.VarcharArray)(u).EncodeText(ci, buf)
}

// EncodeBinary implements pgx interfaces
func (u *ExtendedStringArray) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	return (*pgtype.VarcharArray)(u).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (u *ExtendedStringArray) Scan(src interface{}) error {
	return (*pgtype.VarcharArray)(u).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (u *ExtendedStringArray) Value() (driver.Value, error) {
	return (*pgtype.VarcharArray)(u).Value()
}
