package types

import (
	"strings"

	core "github.com/uhppoted/uhppote-core/types"
)

type Date core.Date

func DateFromString(s string) (*Date, error) {
	date, err := core.DateFromString(s)
	if err != nil {
		return nil, err
	}

	return (*Date)(date), nil
}

func (d *Date) String() string {
	if d != nil {
		return (*core.Date)(d).String()
	}

	return ""
}

func normalise(v string) string {
	return strings.ToLower(strings.ReplaceAll(v, " ", ""))
}
