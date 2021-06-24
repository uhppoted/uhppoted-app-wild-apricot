package types

import (
	"fmt"
	"math"
	"sort"

	"github.com/uhppoted/uhppoted-app-wild-apricot/wild-apricot"
	api "github.com/uhppoted/uhppoted-lib/acl"
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

func MakeGroupList(memberGroups []wildapricot.MemberGroup, displayOrder []string) (*Groups, error) {
	groups := Groups{}
	for _, g := range memberGroups {
		index := uint32(math.MaxUint32)
		for i := range displayOrder {
			name := normalise(displayOrder[i])
			if normalise(g.Name) == name {
				index = uint32(i + 1)
				break
			}
		}

		groups = append(groups, Group{
			ID:    g.ID,
			Name:  g.Name,
			index: index,
		})
	}

	sort.SliceStable(groups, func(i, j int) bool { return groups[i].ID < groups[j].ID })

	return &groups, nil
}

func (groups *Groups) AsTable() *api.Table {
	header := []string{
		"ID",
		"Groups",
	}

	data := [][]string{}

	if groups != nil {
		list := []Group(*groups)

		sort.SliceStable(list, func(i, j int) bool { return normalise(list[i].Name) < normalise(list[j].Name) })
		sort.SliceStable(list, func(i, j int) bool { return list[i].index < list[j].index })

		for _, g := range list {
			row := []string{
				fmt.Sprintf("%v", g.ID),
				fmt.Sprintf("%v", g.Name),
			}

			data = append(data, row)
		}
	}

	table := api.Table{
		Header:  header,
		Records: data,
	}

	return &table
}
