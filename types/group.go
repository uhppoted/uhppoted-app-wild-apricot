package types

import (
	"math"

	"github.com/uhppoted/uhppoted-app-wild-apricot/wild-apricot"
)

type Groups []Group

type Group struct {
	ID    uint32
	Name  string
	index uint32
}

func NewGroup(g wildapricot.MemberGroup) (Group, error) {
	return Group{
		ID:    g.ID,
		Name:  g.Name,
		index: math.MaxUint32,
	}, nil
}
