package types

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/uhppoted/uhppoted-app-wild-apricot/wild-apricot"
)

type Members struct {
	Groups  []Group
	Members []Member
}

type Member struct {
	ID         uint32
	Name       string
	CardNumber *CardNumber
	Active     bool
	Suspended  bool
	Registered *Date
	Expires    *Date
	Groups     map[uint32]struct{}
}

type CardNumber uint32

func (c *CardNumber) String() string {
	if c != nil {
		return fmt.Sprintf("%v", *c)
	}

	return ""
}

func (m *Member) HasRegistered() bool {
	return m != nil && m.Registered != nil
}

func (m *Member) HasExpires() bool {
	return m != nil && m.Expires != nil
}

func MakeMemberList(contacts []wildapricot.Contact, memberGroups []wildapricot.MemberGroup) (*Members, error) {
	groups := []Group{}
	for _, g := range memberGroups {
		groups = append(groups, Group{
			ID:   g.ID,
			Name: g.Name,
		})
	}

	members := []Member{}
	for _, c := range contacts {
		if m, err := transcode(c); err != nil {
			return nil, err
		} else if m != nil {
			members = append(members, *m)
		}
	}
	return &Members{
		Members: members,
		Groups:  groups,
	}, nil
}

func (members *Members) MarshalText() ([]byte, error) {
	return members.MarshalTextIndent("")
}

func (members *Members) MarshalTextIndent(indent string) ([]byte, error) {
	table := [][]string{}

	if members != nil {
		header := []string{
			"ID",
			"Name",
			"Card Number",
			"Active",
			"Suspended",
			"Registered",
			"Expires",
		}

		sort.SliceStable(members.Groups, func(i, j int) bool { return members.Groups[i].ID < members.Groups[j].ID })
		for _, group := range members.Groups {
			header = append(header, group.Name)
		}

		table = append(table, header)

		sort.SliceStable(members.Members, func(i, j int) bool {
			return strings.ToLower(members.Members[i].Name) < strings.ToLower(members.Members[j].Name)
		})

		for _, m := range members.Members {
			row := []string{}
			row = append(row, fmt.Sprintf("%v", m.ID))
			row = append(row, fmt.Sprintf("%v", m.Name))
			row = append(row, fmt.Sprintf("%v", m.CardNumber))
			row = append(row, fmt.Sprintf("%v", m.Active))
			row = append(row, fmt.Sprintf("%v", m.Suspended))
			row = append(row, fmt.Sprintf("%v", m.Registered))
			row = append(row, fmt.Sprintf("%v", m.Expires))

			for _, g := range members.Groups {
				if _, ok := m.Groups[g.ID]; ok {
					row = append(row, "Y")
				} else {
					row = append(row, "N")
				}
			}

			table = append(table, row)
		}
	}

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

func transcode(contact wildapricot.Contact) (*Member, error) {
	member := Member{
		ID:     contact.ID,
		Name:   fmt.Sprintf("%[1]s %[2]s", contact.FirstName, contact.LastName),
		Active: contact.Enabled && strings.ToLower(contact.Status) == "active",
		Groups: map[uint32]struct{}{},
	}

	for _, f := range contact.Fields {
		switch {
		case normalise(f.SystemCode) == "issuspendedmember":
			if v, ok := f.Value.(bool); ok {
				member.Suspended = v
			}

		case normalise(f.SystemCode) == "membersince":
			if v, ok := f.Value.(string); ok {
				if d, err := time.Parse("2006-01-02T15:04:05-07:00", v); err != nil {
					return nil, err
				} else {
					member.Registered = (*Date)(&d)
				}
			}

		case normalise(f.SystemCode) == "renewaldue":
			if v, ok := f.Value.(string); ok {
				if d, err := time.Parse("2006-01-02T15:04:05", v); err != nil {
					return nil, err
				} else {
					expires := d.AddDate(0, 0, -1)
					member.Expires = (*Date)(&expires)
				}
			}

		case normalise(f.Name) == "cardnumber":
			if v, ok := f.Value.(string); ok {
				if n, err := strconv.ParseUint(v, 10, 32); err != nil {
					return nil, err
				} else {
					nn := uint32(n)
					member.CardNumber = (*CardNumber)(&nn)
				}
			}

		case normalise(f.SystemCode) == "groups":
			if groups, ok := f.Value.([]interface{}); ok {
				for _, g := range groups {
					if group, ok := g.(map[string]interface{}); ok {
						if v, ok := group["Id"]; ok {
							if gid, ok := v.(float64); ok {
								member.Groups[uint32(gid)] = struct{}{}
							}
						}
					}
				}
			}
		}
	}

	return &member, nil
}
