package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"

	"github.com/jackc/pgx/pgtype"
)

// ExtendedText is a postgres text field
type ExtendedText pgtype.Text

// MarshalJSON implements json.Marshaler
func (u *ExtendedText) MarshalJSON() ([]byte, error) {
	if u == nil || u.Status == pgtype.Null {
		return []byte("null"), nil
	}

	return json.Marshal(u.String)
}

// StringToText converts a string containing a string into an ExtendedText
func StringToText(input string) ExtendedText {
	return ExtendedText{
		Status: pgtype.Present,
		String: input,
	}
}

// UnmarshalJSON implements json.Unmarshaler
func (u *ExtendedText) UnmarshalJSON(input []byte) error {
	if bytes.Equal(input, []byte("null")) {
		*u = ExtendedText{
			Status: pgtype.Null,
		}
		return nil
	}

	var str string
	if err := json.Unmarshal(input, &str); err != nil {
		return err
	}

	*u = ExtendedText{
		Status: pgtype.Present,
		String: str,
	}

	return nil
}

// Set implements pgx interfaces
func (u *ExtendedText) Set(src interface{}) error {
	return (*pgtype.Text)(u).Set(src)
}

// Get implements pgx interfaces
func (u *ExtendedText) Get() interface{} {
	return (*pgtype.Text)(u).Get()
}

// AssignTo implements pgx interfaces
func (u *ExtendedText) AssignTo(src interface{}) error {
	return (*pgtype.Text)(u).AssignTo(src)
}

// DecodeText implements pgx interfaces
func (u *ExtendedText) DecodeText(ci *pgtype.ConnInfo, src []byte) error {
	return (*pgtype.Text)(u).DecodeText(ci, src)
}

// DecodeBinary implements pgx interfaces
func (u *ExtendedText) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	return (*pgtype.Text)(u).DecodeBinary(ci, src)
}

// EncodeText implements pgx interfaces
func (u *ExtendedText) EncodeText(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	return (*pgtype.Text)(u).EncodeText(ci, buf)
}

// EncodeBinary implements pgx interfaces
func (u *ExtendedText) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	return (*pgtype.Text)(u).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (u *ExtendedText) Scan(src interface{}) error {
	return (*pgtype.Text)(u).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (u *ExtendedText) Value() (driver.Value, error) {
	return (*pgtype.Text)(u).Value()
}
