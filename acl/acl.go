package acl

import (
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	api "github.com/uhppoted/uhppoted-api/acl"
)

type ACL struct {
	doors   []string
	records []record
}

func (acl *ACL) Updated(hash string) bool {
	if hash != "" && hash == acl.Hash() {
		return false
	}

	return true
}

func (acl *ACL) Hash() string {
	if acl != nil {
		header, data := acl.asTable()

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

func (acl *ACL) AsTable() *api.Table {
	header, data := acl.asTable()

	return &api.Table{
		Header:  header,
		Records: data,
	}
}

func (acl *ACL) ToTSV(f io.Writer) error {
	header, data := acl.asTable()

	w := csv.NewWriter(f)
	w.Comma = '\t'

	w.Write(header)
	for _, row := range data {
		w.Write(row)
	}

	w.Flush()

	return nil
}

func (acl *ACL) asTable() ([]string, [][]string) {
	header := []string{}
	data := [][]string{}

	if acl != nil {
		header = append(header, []string{
			"Card Number",
			"From",
			"To",
		}...)

		for _, door := range acl.doors {
			header = append(header, door)
		}

		sort.SliceStable(acl.records, func(i, j int) bool { return acl.records[i].CardNumber < acl.records[j].CardNumber })

		for _, r := range acl.records {
			row := []string{
				fmt.Sprintf("%v", r.CardNumber),
				r.StartDate.Format("2006-01-02"),
				r.EndDate.Format("2006-01-02"),
			}

			for _, door := range acl.doors {
				granted := false
				revoked := false
				profile := -1
				d := normalise(door)

				if _, ok := r.Granted["*"]; ok {
					granted = true
				}

				for k, v := range r.Granted {
					if d == normalise(k) {
						switch vv := v.(type) {
						case bool:
							if vv == true {
								granted = true
							}

						case int:
							if vv >= 2 && vv <= 254 {
								granted = true
								profile = vv
							}

						}
					}
				}

				if _, ok := r.Revoked["*"]; ok {
					revoked = true
				}

				for k, _ := range r.Revoked {
					if d == normalise(k) {
						revoked = true
					}
				}

				if granted && !revoked {
					if profile != -1 {
						row = append(row, fmt.Sprintf("%v", profile))
					} else {
						row = append(row, "Y")
					}
				} else {
					row = append(row, "N")
				}
			}

			data = append(data, row)
		}
	}

	return header, data
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
