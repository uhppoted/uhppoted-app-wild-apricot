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
	Groups  map[uint32]string
	Members []Member
}

type Member struct {
	ID         uint32
	Name       string
	CardNumber *card
	Active     bool
	Suspended  bool
	Registered *date
	Expires    *date
	Groups     map[uint32]struct{}
}

type card uint32

func (c *card) String() string {
	if c != nil {
		return fmt.Sprintf("%v", *c)
	}

	return ""
}

type date time.Time

func (d *date) String() string {
	if d != nil {
		return time.Time(*d).Format("2006-01-02")
	}

	return ""
}

func MakeMemberList(contacts []wildapricot.Contact, memberGroups []wildapricot.MemberGroup) (*Members, error) {
	groups := map[uint32]string{}
	for _, g := range memberGroups {
		groups[g.ID] = g.Name
	}

	members, err := transcode(contacts)
	if err != nil {
		return nil, err
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

		groups := []uint32{}
		for k, _ := range members.Groups {
			groups = append(groups, k)
		}

		sort.SliceStable(groups, func(i, j int) bool { return groups[i] < groups[j] })
		for _, gid := range groups {
			header = append(header, members.Groups[gid])
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

			for _, gid := range groups {
				if _, ok := m.Groups[gid]; ok {
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

func transcode(contacts []wildapricot.Contact) ([]Member, error) {
	members := []Member{}

	for _, c := range contacts {
		member := Member{
			ID:     c.ID,
			Name:   fmt.Sprintf("%[1]s %[2]s", c.FirstName, c.LastName),
			Active: c.Enabled && strings.ToLower(c.Status) == "active",
			Groups: map[uint32]struct{}{},
		}

		for _, f := range c.Fields {
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
						member.Registered = (*date)(&d)
					}
				}

			case normalise(f.SystemCode) == "renewaldue":
				if v, ok := f.Value.(string); ok {
					if d, err := time.Parse("2006-01-02T15:04:05", v); err != nil {
						return nil, err
					} else {
						expires := d.AddDate(0, 0, -1)
						member.Expires = (*date)(&expires)
					}
				}

			case normalise(f.Name) == "cardnumber":
				if v, ok := f.Value.(string); ok {
					if n, err := strconv.ParseUint(v, 10, 32); err != nil {
						return nil, err
					} else {
						nn := uint32(n)
						member.CardNumber = (*card)(&nn)
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

		members = append(members, member)
	}

	return members, nil
}
