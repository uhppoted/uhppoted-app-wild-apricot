package types

import (
	"github.com/uhppoted/uhppoted-app-wild-apricot/wild-apricot"
)

type Group struct {
	ID   uint32
	Name string
}

func NewGroup(g wildapricot.MemberGroup) Group {
	return Group{
		ID:   g.ID,
		Name: g.Name,
	}
}
