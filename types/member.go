package types

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	core "github.com/uhppoted/uhppote-core/types"
	lib "github.com/uhppoted/uhppoted-lib/acl"

	"github.com/uhppoted/uhppoted-app-wild-apricot/log"
	"github.com/uhppoted/uhppoted-app-wild-apricot/wild-apricot"
)

type Members struct {
	Groups  []Group
	Members []Member
}

type MembersWithPIN struct {
	Members
}

type Member struct {
	id         uint32
	Name       string
	CardNumber *CardNumber
	PIN        uint32
	Active     bool
	Suspended  bool
	Registered core.Date
	Expires    core.Date
	Groups     map[uint32]Group
	Membership Membership
	Fields     []Field
}

type CardNumber uint32

func (c *CardNumber) String() string {
	if c != nil {
		return fmt.Sprintf("%v", *c)
	}

	return ""
}

type Membership struct {
	ID   uint32
	Name string
}

type Field struct {
	ID    string
	Name  string
	Value interface{}
}

type field int

const (
	fCardNumber field = iota
	fRegistered
	fExpires
	fSuspended
	fPIN
)

func (f field) String() string {
	return [...]string{"Card Number", "Registered", "Expires", "Suspended"}[f]
}

func (m *Member) Is(membership interface{}) bool {
	if m != nil {
		switch v := membership.(type) {
		case int64:
			if v == int64(m.Membership.ID) {
				return true
			}

		case string:
			if normalise(v) == normalise(m.Membership.Name) {
				return true
			}
		}
	}

	return false
}

func (m *Member) HasCardNumber(card any) bool {
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

func (m *Member) HasPIN() bool {
	if m != nil && m.PIN > 0 && m.PIN < 1000000 {
		return true
	}

	return false
}

func (m *Member) HasRegistered() bool {
	return m != nil && !m.Registered.IsZero()
}

func (m *Member) HasExpires() bool {
	return m != nil && !m.Expires.IsZero()
}

func (m *Member) IsActive() bool {
	return m != nil && m.Active
}

func (m *Member) IsSuspended() bool {
	return m != nil && m.Suspended
}

func (m *Member) HasGroup(group any) bool {
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

func (m *Member) Get(field interface{}) string {
	if m != nil {
		switch v := field.(type) {
		case string:
			vv := normalise(v)
			for _, f := range m.Fields {
				if vv == normalise(f.ID) || vv == normalise(f.Name) {
					return fmt.Sprintf("%v", f.Value)
				}
			}
		}
	}

	return ""
}

func MakeMemberList(contacts []wildapricot.Contact, memberGroups []wildapricot.MemberGroup, cardnumber, pin, facilityCode string, displayOrder []string) (*Members, []error) {
	errors := []error{}

	warnings := struct {
		info     []string
		warnings []string
	}{
		info:     []string{},
		warnings: []string{},
	}

	fields := map[field]string{
		fCardNumber: normalise(cardnumber),
		fPIN:        normalise(pin),
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
		if m, err := transcode(c, groups, fields); err != nil {
			errors = append(errors, fmt.Errorf("Member ID: %d, %v", c.ID, err))
		} else if m != nil {
			if m.CardNumber != nil && *m.CardNumber > 0 && *m.CardNumber < 100000 && facilityCode != "" {
				cardNo := fmt.Sprintf("%v%05v", facilityCode, m.CardNumber)
				if v, err := strconv.ParseUint(cardNo, 10, 32); err != nil {
					warnings.warnings = append(warnings.warnings, fmt.Sprintf("prepending facility code '%v' to card number '%v' for member %v (%v)", facilityCode, m.CardNumber, m.id, err))
				} else {
					warnings.info = append(warnings.info, fmt.Sprintf("prepending facility code '%v' to card %v for member %v\n", facilityCode, m.CardNumber, m.id))
					nn := uint32(v)
					m.CardNumber = (*CardNumber)(&nn)

				}
			}

			members = append(members, *m)
		}
	}

	if len(warnings.warnings) > 0 {
		log.Warnf("error prepending facility code to %v cards", len(warnings.warnings))
	}

	if len(warnings.info) > 0 {
		log.Warnf("prepended facility code to %v cards", len(warnings.info))
	}

	for _, w := range warnings.warnings {
		log.Debugf("%v", w)
	}

	for _, w := range warnings.info {
		log.Debugf("%v", w)
	}

	return &Members{
		Members: members,
		Groups:  groups,
	}, errors
}

func (members *Members) Updated(hash string, withPIN bool) bool {
	if hash != "" && hash == members.hash(withPIN) {
		return false
	}

	return true
}

func (members *Members) Hash() string {
	return members.hash(false)
}

func (members *MembersWithPIN) Hash() string {
	return members.hash(true)
}

func (members *Members) hash(withPIN bool) string {
	asTable := func() ([]string, [][]string) {
		if withPIN {
			return members.asTableWithPIN()
		} else {
			return members.asTable()
		}
	}

	if members != nil {
		header, data := asTable()

		hash := sha256.New()

		for _, h := range header {
			hash.Write([]byte(h))
		}

		for _, r := range data {
			for _, f := range r {
				hash.Write([]byte(f))
			}
		}

		return hex.EncodeToString(hash.Sum(nil))
	}

	return ""
}

func (members *Members) AsTable() *lib.Table {
	header, data := members.asTable()

	return &lib.Table{
		Header:  header,
		Records: data,
	}
}

func (members *Members) AsTableWithPIN() *lib.Table {
	header, data := members.asTableWithPIN()

	return &lib.Table{
		Header:  header,
		Records: data,
	}
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

func (members *Members) ToTSVWithPIN(f io.Writer) error {
	header, data := members.asTableWithPIN()

	w := csv.NewWriter(f)
	w.Comma = '\t'

	w.Write(header)
	for _, row := range data {
		w.Write(row)
	}

	w.Flush()

	return nil
}

func (members *Members) asTable() ([]string, [][]string) {
	header := []string{
		"Name",
		"Card Number",
		"Membership",
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
			row = append(row, fmt.Sprintf("%v", m.Membership.Name))
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

func (members *Members) asTableWithPIN() ([]string, [][]string) {
	header := []string{
		"Name",
		"Card Number",
		"PIN",
		"Membership",
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
			var pin string

			if m.PIN != 0 {
				pin = fmt.Sprintf("%v", m.PIN)
			} else {
				pin = ""
			}

			row := []string{}
			row = append(row, fmt.Sprintf("%v", m.Name))
			row = append(row, fmt.Sprintf("%v", m.CardNumber))
			row = append(row, fmt.Sprintf("%v", pin))
			row = append(row, fmt.Sprintf("%v", m.Membership.Name))
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

func transcode(contact wildapricot.Contact, sysgroups []Group, fields map[field]string) (*Member, error) {
	member := Member{
		id:   contact.ID,
		Name: fmt.Sprintf("%[1]s %[2]s", contact.FirstName, contact.LastName),
		Membership: Membership{
			ID:   contact.MembershipLevel.ID,
			Name: contact.MembershipLevel.Name,
		},
		Active: contact.Enabled && strings.ToLower(contact.Status) == "active",
		Groups: map[uint32]Group{},
		Fields: []Field{},
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
					return nil, fmt.Errorf("unable to parse 'Member since' date '%v' (%v)", v, err)
				} else {
					member.Registered = core.Date(d)
				}
			}

		case normalise(f.SystemCode) == fields[fExpires]:
			if v, ok := f.Value.(string); ok {
				if d, err := time.Parse("2006-01-02T15:04:05", v); err != nil {
					return nil, fmt.Errorf("unable to parse 'Renewal' date '%v' (%v)", v, err)
				} else {
					expires := d.AddDate(0, 0, -1)
					member.Expires = core.Date(expires)
				}
			}

		case normalise(f.Name) == fields[fCardNumber]:
			if v, ok := f.Value.(string); ok {
				if v != "" {
					if n, err := strconv.ParseUint(v, 10, 32); err != nil {
						return nil, fmt.Errorf("error parsing card number '%v' (%v)", v, err)
					} else {
						nn := uint32(n)
						member.CardNumber = (*CardNumber)(&nn)
					}
				}
			}

		case normalise(f.Name) == fields[fPIN]:
			if v, ok := f.Value.(string); ok {
				if v != "" {
					if n, err := strconv.ParseUint(v, 10, 32); err != nil {
						return nil, fmt.Errorf("error parsing PIN '%v' (%v)", v, err)
					} else {
						member.PIN = uint32(n)
					}
				}
			}

		case normalise(f.SystemCode) == "groups":
			if groups, ok := f.Value.([]any); ok {
				for _, g := range groups {
					if gg, ok := g.(map[string]any); ok {
						group := Group{}

						if v, ok := gg["Id"]; ok {
							if id, ok := v.(float64); ok {
								group.ID = uint32(id)

								for _, sysgroup := range sysgroups {
									if group.ID == sysgroup.ID {
										group.Name = sysgroup.Name
									}
								}
							}
						}

						// Ref. https://github.com/uhppoted/uhppoted-app-wild-apricot/issues/6
						//      Replaced with actual group name from groups list - Wild Apricot does not keep
						//      the member group names synchronized with the actual group names leading to
						//      confusion when resolving permissions by group name.
						// if v, ok := gg["Label"]; ok {
						// 	if name, ok := v.(string); ok {
						// 		group.Name = name
						// 	}
						// }

						// ... warn if member group name does not match group name (for legacy systems)
						if v, ok := gg["Label"]; ok {
							if name, ok := v.(string); ok {
								if name != group.Name {
									log.Warnf("%v: group name (%v) does not match groups name (%v)", member.Name, name, group.Name)
								}
							}
						}

						member.Groups[group.ID] = group
					}
				}
			}
		}
	}

	for _, f := range contact.Fields {
		member.Fields = append(member.Fields, Field{
			ID:    f.SystemCode,
			Name:  f.Name,
			Value: f.Value,
		})

	}

	return &member, nil
}
