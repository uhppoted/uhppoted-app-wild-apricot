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
