package types

import (
	"regexp"
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
	re := regexp.MustCompile(`[^a-z1-9]`)

	return re.ReplaceAllString(strings.ToLower(v), "")
}
