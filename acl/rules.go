package acl

import (
	"bytes"
	"sort"
	"time"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"

	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
)

type ACL []record

type Rules struct {
	library *ast.KnowledgeLibrary
}

func NewRules(file string) (*Rules, error) {
	kb := ast.NewKnowledgeLibrary()
	if err := builder.NewRuleBuilder(kb).BuildRuleFromResource("acl", "0.0.0", pkg.NewFileResource(file)); err != nil {
		return nil, err
	}

	return &Rules{
		library: kb,
	}, nil
}

func (rules *Rules) MakeACL(members types.Members) (ACL, error) {
	acl := ACL{}

	startDate := startOfYear()

	for _, m := range members.Members {
		if m.CardNumber != nil {
			r := record{
				ID:         m.ID,
				Name:       m.Name,
				CardNumber: uint32(*m.CardNumber),
				StartDate:  startDate,
			}

			if err := rules.eval(m, &r); err != nil {
				return nil, err
			}

			acl = append(acl, r)
		}
	}

	sort.SliceStable(acl, func(i, j int) bool { return acl[i].ID < acl[j].ID })

	return acl, nil
}

func (rules *Rules) eval(m types.Member, r *record) error {
	context := ast.NewDataContext()

	if err := context.Add("member", &m); err != nil {
		return err
	}

	if err := context.Add("record", r); err != nil {
		return err
	}

	enjin := engine.NewGruleEngine()
	kb := rules.library.NewKnowledgeBaseInstance("acl", "0.0.0")

	return enjin.Execute(context, kb)
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

func startOfYear() time.Time {
	return time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.Local)
}
