package acl

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"
)

type ACL struct {
	doors   []string
	records []record
}

func (acl *ACL) MarshalTextIndent(indent string) ([]byte, error) {
	table := [][]string{}

	if acl != nil {
		header := []string{
			"ID",
			"Name",
			"Card Number",
			"Start Date",
			"End Date",
		}

		sort.SliceStable(acl.doors, func(i, j int) bool { return acl.doors[i] < acl.doors[j] })
		for _, door := range acl.doors {
			header = append(header, door)
		}

		table = append(table, header)

		sort.SliceStable(acl.records, func(i, j int) bool { return acl.records[i].CardNumber < acl.records[j].CardNumber })

		for _, r := range acl.records {
			row := []string{}
			row = append(row, fmt.Sprintf("%v", r.ID))
			row = append(row, fmt.Sprintf("%v", r.Name))
			row = append(row, fmt.Sprintf("%v", r.CardNumber))
			row = append(row, r.StartDate.Format("2006-01-02"))
			row = append(row, r.EndDate.Format("2006-01-02"))

			for _, door := range acl.doors {
				granted := false
				revoked := false

				if _, ok := r.Granted["*"]; ok {
					granted = true
				}

				if _, ok := r.Granted[normalise(door)]; ok {
					granted = true
				}

				if _, ok := r.Revoked["*"]; ok {
					revoked = true
				}

				if _, ok := r.Revoked[normalise(door)]; ok {
					revoked = true
				}

				if granted && !revoked {
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

func normalise(v string) string {
	return strings.ToLower(strings.ReplaceAll(v, " ", ""))
}

func startOfYear() time.Time {
	return time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.Local)
}

func endOfYear() time.Time {
	return time.Date(time.Now().Year(), time.December, 31, 23, 59, 59, 0, time.Local)
}
