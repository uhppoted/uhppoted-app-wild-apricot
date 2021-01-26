package rules

import (
	"bytes"
	"sort"

	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
)

type Rules {
}

type ACL []record

type record struct {
	ID         uint32
	Name       string
	CardNumber uint32
}

func MakeACL(members types.Members) (ACL, error) {
	acl := ACL{}

	for _, m := range members.Members {
		if m.CardNumber != nil {
			r := record{
				ID:         m.ID,
				Name:       m.Name,
				CardNumber: uint32(*m.CardNumber),
			}

			acl = append(acl, r)
		}
	}

	sort.SliceStable(acl, func(i, j int) bool { return acl[i].ID < acl[j].ID })

	return acl, nil
}

func (a *ACL) MarshalText() ([]byte, error) {
	return a.MarshalTextIndent("")
}

func (a *ACL) MarshalTextIndent(indent string) ([]byte, error) {
	//	table := [][]string{}
	//
	//	if members != nil {
	//		header := []string{
	//			"ID",
	//			"Name",
	//			"Card Number",
	//			"Active",
	//			"Suspended",
	//			"Registered",
	//			"Expires",
	//		}
	//
	//		sort.SliceStable(members.Groups, func(i, j int) bool { return members.Groups[i].ID < members.Groups[j].ID })
	//		for _, group := range members.Groups {
	//			header = append(header, group.Name)
	//		}
	//
	//		table = append(table, header)
	//
	//		sort.SliceStable(members.Members, func(i, j int) bool {
	//			return strings.ToLower(members.Members[i].Name) < strings.ToLower(members.Members[j].Name)
	//		})
	//
	//		for _, m := range members.Members {
	//			row := []string{}
	//			row = append(row, fmt.Sprintf("%v", m.ID))
	//			row = append(row, fmt.Sprintf("%v", m.Name))
	//			row = append(row, fmt.Sprintf("%v", m.CardNumber))
	//			row = append(row, fmt.Sprintf("%v", m.Active))
	//			row = append(row, fmt.Sprintf("%v", m.Suspended))
	//			row = append(row, fmt.Sprintf("%v", m.Registered))
	//			row = append(row, fmt.Sprintf("%v", m.Expires))
	//
	//			for _, g := range members.Groups {
	//				if _, ok := m.Groups[g.ID]; ok {
	//					row = append(row, "Y")
	//				} else {
	//					row = append(row, "N")
	//				}
	//			}
	//
	//			table = append(table, row)
	//		}
	//	}
	//
	var b bytes.Buffer

	//	if len(table) > 0 {
	//		widths := make([]int, len(table[0]))
	//		for _, row := range table {
	//			for i, field := range row {
	//				if len(field) > widths[i] {
	//					widths[i] = len(field)
	//				}
	//			}
	//		}
	//
	//		for i := 1; i < len(widths); i++ {
	//			widths[i-1] += 1
	//		}
	//
	//		for _, row := range table {
	//			fmt.Fprintf(&b, "%s", indent)
	//			for i, field := range row {
	//				fmt.Fprintf(&b, "%-*v", widths[i], field)
	//			}
	//			fmt.Fprintln(&b)
	//		}
	//	}

	return b.Bytes(), nil
}
