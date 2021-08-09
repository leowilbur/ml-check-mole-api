package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/pgtype"
)

// ExtendedUUID is Postgres UUID with JSON encoding support
type ExtendedUUID pgtype.UUID

// parseUUID converts a string UUID in standard form to a byte array.
func parseUUID(src string) (dst [16]byte, err error) {
	src = src[0:8] + src[9:13] + src[14:18] + src[19:23] + src[24:]
	buf, err := hex.DecodeString(src)
	if err != nil {
		return dst, err
	}

	copy(dst[:], buf)
	return dst, err
}

// encodeUUID converts a uuid byte array to UUID standard string form.
func encodeUUID(src [16]byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", src[0:4], src[4:6], src[6:8], src[8:10], src[10:16])
}

// String implements Stringer
func (u *ExtendedUUID) String() string {
	if u == nil || u.Status != pgtype.Present {
		return ""
	}
	return encodeUUID(u.Bytes)
}

// MarshalJSON implements json.Marshaler
func (u *ExtendedUUID) MarshalJSON() ([]byte, error) {
	if u == nil || u.Status == pgtype.Null {
		return []byte("null"), nil
	}

	return []byte(`"` + encodeUUID(u.Bytes) + `"`), nil
}

// StringToUUID converts a string containing a properly formatted UUID into
// an ExtendedUUID
func StringToUUID(input string) (ExtendedUUID, error) {
	parsed, err := parseUUID(input)
	if err != nil {
		return ExtendedUUID{}, err
	}

	return ExtendedUUID{
		Status: pgtype.Present,
		Bytes:  parsed,
	}, nil
}

// UnmarshalJSON implements json.Unmarshaler
func (u *ExtendedUUID) UnmarshalJSON(input []byte) error {
	if bytes.Equal(input, []byte("null")) {
		*u = ExtendedUUID{
			Status: pgtype.Null,
		}
		return nil
	}

	var str string
	if err := json.Unmarshal(input, &str); err != nil {
		return err
	}

	uuid, err := StringToUUID(str)
	if err != nil {
		return err
	}

	*u = uuid

	return nil
}

// Set implements pgx interfaces
func (u *ExtendedUUID) Set(src interface{}) error {
	return (*pgtype.UUID)(u).Set(src)
}

// Get implements pgx interfaces
func (u *ExtendedUUID) Get() interface{} {
	return (*pgtype.UUID)(u).Get()
}

// AssignTo implements pgx interfaces
func (u *ExtendedUUID) AssignTo(src interface{}) error {
	return (*pgtype.UUID)(u).AssignTo(src)
}

// DecodeText implements pgx interfaces
func (u *ExtendedUUID) DecodeText(ci *pgtype.ConnInfo, src []byte) error {
	return (*pgtype.UUID)(u).DecodeText(ci, src)
}

// DecodeBinary implements pgx interfaces
func (u *ExtendedUUID) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	return (*pgtype.UUID)(u).DecodeBinary(ci, src)
}

// EncodeText implements pgx interfaces
func (u *ExtendedUUID) EncodeText(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	return (*pgtype.UUID)(u).EncodeText(ci, buf)
}

// EncodeBinary implements pgx interfaces
func (u *ExtendedUUID) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	return (*pgtype.UUID)(u).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (u *ExtendedUUID) Scan(src interface{}) error {
	return (*pgtype.UUID)(u).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (u *ExtendedUUID) Value() (driver.Value, error) {
	return (*pgtype.UUID)(u).Value()
}
