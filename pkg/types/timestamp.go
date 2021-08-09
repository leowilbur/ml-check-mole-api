package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/pgtype"
)

// ExtendedTimestamp is a postgres timestamp field
type ExtendedTimestamp pgtype.Timestamp

// MarshalJSON implements json.Marshaler
func (u *ExtendedTimestamp) MarshalJSON() ([]byte, error) {
	if u == nil || u.Status == pgtype.Null {
		return []byte("null"), nil
	}

	return json.Marshal(u.Time.Unix())
}

// TimeToTimestamp converts a time.Time into an ExtendedTimestamp
func TimeToTimestamp(input time.Time) ExtendedTimestamp {
	return ExtendedTimestamp{
		Status: pgtype.Present,
		Time:   input,
	}
}

// UnmarshalJSON implements json.Unmarshaler
func (u *ExtendedTimestamp) UnmarshalJSON(input []byte) error {
	if bytes.Equal(input, []byte("null")) {
		*u = ExtendedTimestamp{
			Status: pgtype.Null,
		}
		return nil
	}

	var str interface{}
	if err := json.Unmarshal(input, &str); err != nil {
		return err
	}

	switch v := str.(type) {
	case float64:
		*u = ExtendedTimestamp{
			Status: pgtype.Present,
			Time:   time.Unix(int64(v), 0),
		}
		return nil
	case string:
		ts, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
		*u = ExtendedTimestamp{
			Status: pgtype.Present,
			Time:   ts,
		}
		return nil
	default:
		return errors.New("unknown timestamp input type")
	}
}

// Set implements pgx interfaces
func (u *ExtendedTimestamp) Set(src interface{}) error {
	return (*pgtype.Timestamp)(u).Set(src)
}

// Get implements pgx interfaces
func (u *ExtendedTimestamp) Get() interface{} {
	return (*pgtype.Timestamp)(u).Get()
}

// AssignTo implements pgx interfaces
func (u *ExtendedTimestamp) AssignTo(src interface{}) error {
	return (*pgtype.Timestamp)(u).AssignTo(src)
}

// DecodeText implements pgx interfaces
func (u *ExtendedTimestamp) DecodeText(ci *pgtype.ConnInfo, src []byte) error {
	return (*pgtype.Timestamp)(u).DecodeText(ci, src)
}

// DecodeBinary implements pgx interfaces
func (u *ExtendedTimestamp) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	return (*pgtype.Timestamp)(u).DecodeBinary(ci, src)
}

// EncodeText implements pgx interfaces
func (u *ExtendedTimestamp) EncodeText(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	return (*pgtype.Timestamp)(u).EncodeText(ci, buf)
}

// EncodeBinary implements pgx interfaces
func (u *ExtendedTimestamp) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	return (*pgtype.Timestamp)(u).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (u *ExtendedTimestamp) Scan(src interface{}) error {
	return (*pgtype.Timestamp)(u).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (u *ExtendedTimestamp) Value() (driver.Value, error) {
	return (*pgtype.Timestamp)(u).Value()
}
