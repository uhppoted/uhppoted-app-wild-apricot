package acl

import (
	"regexp"
	"strconv"
	"time"

	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
)

type record struct {
	Name       string
	CardNumber uint32
	PIN        uint32
	StartDate  time.Time
	EndDate    time.Time
	Granted    map[string]interface{}
	Revoked    map[string]struct{}
}

func (r *record) SetCardNumber(card interface{}) {
	if r != nil {
		switch v := card.(type) {
		case string:
			if v, err := strconv.ParseUint(v, 10, 32); err == nil {
				r.CardNumber = uint32(v)
			}

		case int64:
			r.CardNumber = uint32(v)
		}
	}
}

func (r *record) SetStartDate(t interface{}) {
	if r != nil {
		switch v := t.(type) {
		case string:
			if date, err := time.Parse("2006-01-02", v); err == nil {
				r.StartDate = date
			}

		case *time.Time:
			if v != nil {
				r.StartDate = *v
			}

		case *types.Date:
			if v != nil {
				r.StartDate = time.Time(*v)
			}
		}
	}
}

func (r *record) SetEndDate(t interface{}) {
	if r != nil {
		switch v := t.(type) {
		case string:
			if date, err := time.Parse("2006-01-02", v); err == nil {
				r.EndDate = date
			}

		case *time.Time:
			if v != nil {
				r.EndDate = *v
			}

		case *types.Date:
			if v != nil {
				r.EndDate = time.Time(*v)
			}
		}
	}
}

func (r *record) Grant(permissions ...interface{}) {
	if r != nil {
		// parse Grant(door, profile)
		if len(permissions) == 2 {
			if door, ok := permissions[0].(string); ok {
				switch profile := permissions[1].(type) {
				case int:
					r.Granted[normalise(door)] = profile
					return

				case int64:
					r.Granted[normalise(door)] = int(profile)
					return
				}
			}
		}

		// parse Grant(door...)
		for _, p := range permissions {
			if d, ok := p.(string); ok {
				if match := regexp.MustCompile(`(\S.*?):([0-9]+)`).FindStringSubmatch(d); match != nil {
					door := normalise(match[1])
					profile, _ := strconv.Atoi(match[2])
					r.Granted[door] = profile
				} else {
					r.Granted[normalise(d)] = true
				}
			}
		}
	}
}

func (r *record) Revoke(door ...string) {
	if r != nil {
		for _, d := range door {
			r.Revoked[normalise(d)] = struct{}{}
		}
	}
}
