package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// dateLayout is the wire/storage format for Date (calendar date, no time).
const dateLayout = "2006-01-02"

// Date is a calendar date stored as a MySQL DATE column and (de)serialized as
// "YYYY-MM-DD" in JSON. The zero value marshals to null / stores as NULL.
type Date struct {
	time.Time
}

// GormDataType maps Date to a DATE column.
func (Date) GormDataType() string { return "date" }

// MarshalJSON renders the date as "YYYY-MM-DD" (or null when zero).
func (d Date) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Time.Format(dateLayout) + `"`), nil
}

// UnmarshalJSON accepts "YYYY-MM-DD" (or empty/null -> zero).
func (d *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return fmt.Errorf("Date: expected YYYY-MM-DD: %w", err)
	}
	d.Time = t
	return nil
}

// Value implements driver.Valuer — a zero date persists as NULL.
func (d Date) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time.Format(dateLayout), nil
}

// Scan implements sql.Scanner for time.Time / string / []byte sources.
func (d *Date) Scan(value any) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		d.Time = v
	case []byte:
		return d.UnmarshalJSON([]byte(`"` + string(v) + `"`))
	case string:
		return d.UnmarshalJSON([]byte(`"` + v + `"`))
	default:
		return fmt.Errorf("Date: unsupported scan type %T", value)
	}
	return nil
}

// JSONRaw is arbitrary JSON stored in a MySQL JSON column. It round-trips the
// raw bytes untouched and (de)serializes transparently in API payloads, so
// callers can read/write nested objects/arrays without a fixed Go shape.
type JSONRaw json.RawMessage

// MarshalJSON emits the stored bytes verbatim (null when empty).
func (j JSONRaw) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return j, nil
}

// UnmarshalJSON copies the incoming bytes verbatim.
func (j *JSONRaw) UnmarshalJSON(b []byte) error {
	if j == nil {
		return fmt.Errorf("JSONRaw: UnmarshalJSON on nil pointer")
	}
	*j = append((*j)[0:0], b...)
	return nil
}

// Value implements driver.Valuer — empty marshals to NULL.
func (j JSONRaw) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

// Scan implements sql.Scanner for []byte / string JSON sources.
func (j *JSONRaw) Scan(value any) error {
	if value == nil {
		*j = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*j = append((*j)[0:0], v...)
	case string:
		*j = append((*j)[0:0], []byte(v)...)
	default:
		return fmt.Errorf("JSONRaw: unsupported scan type %T", value)
	}
	return nil
}

// StringSlice is a []string persisted as a JSON column (MySQL JSON type).
type StringSlice []string

// Value implements driver.Valuer — marshals to a JSON array string.
func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan implements sql.Scanner — unmarshals a JSON array from the DB.
func (s *StringSlice) Scan(value any) error {
	if value == nil {
		*s = StringSlice{}
		return nil
	}
	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("StringSlice: unsupported scan type %T", value)
	}
	if len(b) == 0 {
		*s = StringSlice{}
		return nil
	}
	return json.Unmarshal(b, s)
}
