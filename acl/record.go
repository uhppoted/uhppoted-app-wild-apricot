package acl

import (
	"time"

	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
)

type record struct {
	ID         uint32
	Name       string
	CardNumber uint32
	StartDate  time.Time
	EndDate    time.Time
}

func (r *record) SetStartDate(t interface{}) {
	if r == nil {
		return
	}

	switch v := t.(type) {
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

func (r *record) SetEndDate(t interface{}) {
	if r == nil {
		return
	}

	switch v := t.(type) {
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
