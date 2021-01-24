package commands

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/uhppoted/uhppoted-app-wild-apricot/wild-apricot"
)

type Members struct {
	Members []Member
}

type Member struct {
	ID         uint32
	Name       string
	Active     bool
	Suspended  bool
	Registered *date
	Renew      *date
}

type date time.Time

func (d *date) String() string {
	if d != nil {
		return time.Time(*d).Format("2006-01-02")
	}

	return ""
}

func makeMemberList(contacts []wildapricot.Contact, groups []wildapricot.Group) (*Members, error) {
	members := []Member{}

	for _, c := range contacts {
		member := Member{
			ID:         c.ID,
			Name:       c.Name,
			Active:     c.Active,
			Suspended:  c.Suspended,
			Registered: (*date)(c.MemberSince),
			Renew:      (*date)(c.Renew),
		}

		members = append(members, member)
	}

	return &Members{
		Members: members,
	}, nil
}

func (members *Members) MarshalText() ([]byte, error) {
	return members.MarshalTextIndent("")
}

func (members *Members) MarshalTextIndent(indent string) ([]byte, error) {
	var b bytes.Buffer

	if members != nil {
		sort.SliceStable(members.Members, func(i, j int) bool {
			return strings.ToLower(members.Members[i].Name) < strings.ToLower(members.Members[j].Name)
		})

		table := [][]string{}
		for _, m := range members.Members {
			row := []string{}
			row = append(row, fmt.Sprintf("%v", m.ID))
			row = append(row, fmt.Sprintf("%v", m.Name))
			row = append(row, fmt.Sprintf("%v", m.Active))
			row = append(row, fmt.Sprintf("%v", m.Suspended))
			row = append(row, fmt.Sprintf("%v", m.Registered))
			row = append(row, fmt.Sprintf("%v", m.Renew))

			table = append(table, row)
		}

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
	}

	return b.Bytes(), nil
}
