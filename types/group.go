package types

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"sort"

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

func (groups *Groups) MarshalText() ([]byte, error) {
	return groups.MarshalTextIndent("")
}

func (groups *Groups) MarshalTextIndent(indent string) ([]byte, error) {
	header, data := groups.asTable()
	table := [][]string{}

	table = append(table, header)
	table = append(table, data...)

	var b bytes.Buffer

	if len(table) > 0 {
		widths := make([]int, len(table[0]))
		for _, row := range table {
			for i, field := range row {
				if len(field) > widths[i] {
					widths[i] = len(field)
				}
			}
		}

		for i := 1; i < len(widths); i++ {
			widths[i-1] += 1
		}

		for _, row := range table {
			fmt.Fprintf(&b, "%s", indent)
			for i, field := range row {
				fmt.Fprintf(&b, "%-*v", widths[i], field)
			}
			fmt.Fprintln(&b)
		}
	}

	return b.Bytes(), nil
}

func (groups *Groups) ToTSV(f io.Writer) error {
	header, data := groups.asTable()

	w := csv.NewWriter(f)
	w.Comma = '\t'

	w.Write(header)
	for _, row := range data {
		w.Write(row)
	}

	w.Flush()

	return nil
}

func (groups *Groups) asTable() ([]string, [][]string) {
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

	return header, data
}
