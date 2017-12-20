package urn

import (
	"database/sql/driver"
	"fmt"
	"net/url"
	"strings"
)

var (
	// Nil represents a completely empty (invalid) URN
	Nil URN
)

// URN is a Universal Resource Name
type URN struct {
	Domain  string
	Service string
	Name    string
	ID      string
}

// Parse string to URN.
func Parse(s string) (URN, error) {
	uri, err := url.Parse(strings.ToLower(s))
	if err != nil {
		return Nil, err
	}

	u := URN{
		Domain:  uri.Scheme,
		Service: uri.Host,
	}

	paths := strings.Split(uri.Path, "/")
	if len(paths) > 1 {
		u.Name = paths[1]
	}
	if len(paths) > 2 {
		u.ID = paths[2]
	}

	return u, nil
}

// String serializes to a string.
func (u URN) String() string {
	return fmt.Sprintf("%s://%s/%s/%s", strings.ToLower(u.Domain), strings.ToLower(u.Service), strings.ToLower(u.Name), strings.ToLower(u.ID))
}

// New returns a new ID.
func New(domain, service, name, id string) URN {
	return URN{
		Domain:  strings.ToLower(domain),
		Service: strings.ToLower(service),
		Name:    strings.ToLower(name),
		ID:      strings.ToLower(id),
	}
}

// MarshalText implements encoding.TextMarshaler.
func (u URN) MarshalText() ([]byte, error) {
	return []byte(u.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (u *URN) UnmarshalText(b []byte) error {
	urn, err := Parse(string(b))
	if err != nil {
		return err
	}
	*u = urn
	return nil
}

// Value converts the URN into a SQL driver value which can be used to
// directly use the URN as parameter to a SQL query.
func (u URN) Value() (driver.Value, error) {
	if u.IsNil() {
		return nil, nil
	}
	return u.String(), nil
}

// Scan implements the sql.Scanner interface. It supports converting from
// string, []byte, or nil into a URN value. Attempting to convert from
// another type will return an error.
func (u *URN) Scan(src interface{}) error {
	switch v := src.(type) {
	case nil:
		*u = Nil
		return nil
	case []byte:
		return u.UnmarshalText(v)
	case string:
		return u.UnmarshalText([]byte(v))
	default:
		return fmt.Errorf("Scan: unable to scan type %T into URN", v)
	}
}

// IsNil returns true if this is a "nil" URN.
func (u URN) IsNil() bool {
	return u == Nil
}
