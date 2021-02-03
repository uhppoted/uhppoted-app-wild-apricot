package types

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"math"
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
	id         uint32
	Name       string
	CardNumber *CardNumber
	Active     bool
	Suspended  bool
	Registered *Date
	Expires    *Date
	Groups     map[uint32]Group
}

type CardNumber uint32

func (c *CardNumber) String() string {
	if c != nil {
		return fmt.Sprintf("%v", *c)
	}

	return ""
}

type Field int

const (
	fCardNumber Field = iota
	fRegistered
	fExpires
	fSuspended
)

func (f Field) String() string {
	return [...]string{"Card Number", "Registered", "Expires", "Suspended"}[f]
}

func (m *Member) HasCardNumber(card interface{}) bool {
	if m != nil && m.CardNumber != nil {
		switch v := card.(type) {
		case int64:
			if v == int64(*m.CardNumber) {
				return true
			}
		}
	}

	return false
}

func (m *Member) HasRegistered() bool {
	return m != nil && m.Registered != nil
}

func (m *Member) HasExpires() bool {
	return m != nil && m.Expires != nil
}

func (m *Member) IsActive() bool {
	return m != nil && m.Active == true
}

func (m *Member) IsSuspended() bool {
	return m != nil && m.Suspended
}

func (m *Member) HasGroup(group interface{}) bool {
	if m != nil {
		switch v := group.(type) {
		case string:
			vv := normalise(v)
			for _, g := range m.Groups {
				if vv == normalise(g.Name) {
					return true
				}
			}

		case int64:
			for _, g := range m.Groups {
				if v == int64(g.ID) {
					return true
				}
			}
		}
	}

	return false
}

func MakeMemberList(contacts []wildapricot.Contact, memberGroups []wildapricot.MemberGroup, cardnumber string, displayOrder []string) (*Members, error) {
	fields := map[Field]string{
		fCardNumber: normalise(cardnumber),
		fRegistered: normalise("MemberSince"),
		fSuspended:  normalise("IsSuspendedMember"),
		fExpires:    normalise("RenewalDue"),
	}

	groups := []Group{}
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

	members := []Member{}
	for _, c := range contacts {
		if m, err := transcode(c, fields); err != nil {
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

func (members *Members) ToTSV(f io.Writer) error {
	header, data := members.asTable()

	w := csv.NewWriter(f)
	w.Comma = '\t'

	w.Write(header)
	for _, row := range data {
		w.Write(row)
	}

	w.Flush()

	return nil
}

func (members *Members) MarshalText() ([]byte, error) {
	return members.MarshalTextIndent("")
}

func (members *Members) MarshalTextIndent(indent string) ([]byte, error) {
	header, data := members.asTable()
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

func (members *Members) asTable() ([]string, [][]string) {
	header := []string{
		"Name",
		"Card Number",
		"Active",
		"Suspended",
		"Registered",
		"Expires",
	}

	data := [][]string{}

	f := func(b bool) string {
		if b {
			return "Y"
		}

		return "N"
	}

	if members != nil {
		sort.SliceStable(members.Groups, func(i, j int) bool { return normalise(members.Groups[i].Name) < normalise(members.Groups[j].Name) })
		sort.SliceStable(members.Groups, func(i, j int) bool { return members.Groups[i].index < members.Groups[j].index })

		for _, group := range members.Groups {
			header = append(header, group.Name)
		}

		sort.SliceStable(members.Members, func(i, j int) bool {
			return strings.ToLower(members.Members[i].Name) < strings.ToLower(members.Members[j].Name)
		})

		for _, m := range members.Members {
			row := []string{}
			row = append(row, fmt.Sprintf("%v", m.Name))
			row = append(row, fmt.Sprintf("%v", m.CardNumber))
			row = append(row, f(m.Active))
			row = append(row, f(m.Suspended))
			row = append(row, fmt.Sprintf("%v", m.Registered))
			row = append(row, fmt.Sprintf("%v", m.Expires))

			for _, g := range members.Groups {
				if _, ok := m.Groups[g.ID]; ok {
					row = append(row, "Y")
				} else {
					row = append(row, "N")
				}
			}

			data = append(data, row)
		}
	}

	return header, data
}

func transcode(contact wildapricot.Contact, fields map[Field]string) (*Member, error) {
	member := Member{
		id:     contact.ID,
		Name:   fmt.Sprintf("%[1]s %[2]s", contact.FirstName, contact.LastName),
		Active: contact.Enabled && strings.ToLower(contact.Status) == "active",
		Groups: map[uint32]Group{},
	}

	for _, f := range contact.Fields {
		switch {
		case normalise(f.SystemCode) == fields[fSuspended]:
			if v, ok := f.Value.(bool); ok {
				member.Suspended = v
			}

		case normalise(f.SystemCode) == fields[fRegistered]:
			if v, ok := f.Value.(string); ok {
				if d, err := time.Parse("2006-01-02T15:04:05-07:00", v); err != nil {
					return nil, fmt.Errorf("Unable to parse 'Member since' date '%v' (%v)", v, err)
				} else {
					member.Registered = (*Date)(&d)
				}
			}

		case normalise(f.SystemCode) == fields[fExpires]:
			if v, ok := f.Value.(string); ok {
				if d, err := time.Parse("2006-01-02T15:04:05", v); err != nil {
					return nil, fmt.Errorf("Unable to parse 'Renewal' date '%v' (%v)", v, err)
				} else {
					expires := d.AddDate(0, 0, -1)
					member.Expires = (*Date)(&expires)
				}
			}

		case normalise(f.Name) == fields[fCardNumber]:
			if v, ok := f.Value.(string); ok {
				if v != "" {
					if n, err := strconv.ParseUint(v, 10, 32); err != nil {
						return nil, fmt.Errorf("Error parsing card number '%v' (%v)", v, err)
					} else {
						nn := uint32(n)
						member.CardNumber = (*CardNumber)(&nn)
					}
				}
			}

		case normalise(f.SystemCode) == "groups":
			if groups, ok := f.Value.([]interface{}); ok {
				for _, g := range groups {
					if gg, ok := g.(map[string]interface{}); ok {
						group := Group{}

						if v, ok := gg["Id"]; ok {
							if id, ok := v.(float64); ok {
								group.ID = uint32(id)
							}
						}

						if v, ok := gg["Label"]; ok {
							if name, ok := v.(string); ok {
								group.Name = name
							}
						}

						member.Groups[group.ID] = group
					}
				}
			}
		}
	}

	return &member, nil
}
