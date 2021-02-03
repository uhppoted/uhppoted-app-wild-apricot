package acl

import (
	"time"

	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
)

type record struct {
	Name       string
	CardNumber uint32
	StartDate  time.Time
	EndDate    time.Time
	Granted    map[string]struct{}
	Revoked    map[string]struct{}
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

func (r *record) Grant(door string) {
	if r != nil {
		r.Granted[normalise(door)] = struct{}{}
	}
}

func (r *record) Revoke(door string) {
	if r != nil {
		r.Revoked[normalise(door)] = struct{}{}
	}
}
