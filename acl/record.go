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

func (r *record) Grant(door ...string) {
	if r != nil {
		for _, d := range door {
			if match := regexp.MustCompile(`(\S.*?):([0-9]+)`).FindStringSubmatch(d); match != nil {
				permission := normalise(match[1])
				profile, _ := strconv.Atoi(match[2])
				r.Granted[normalise(permission)] = profile
			} else {
				r.Granted[normalise(d)] = true
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
