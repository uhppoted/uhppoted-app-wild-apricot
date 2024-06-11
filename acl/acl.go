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

	core "github.com/uhppoted/uhppote-core/types"
	lib "github.com/uhppoted/uhppoted-lib/acl"
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

func (acl *ACL) AsTable() *lib.Table {
	header, data := acl.asTable()

	return &lib.Table{
		Header:  header,
		Records: data,
	}
}

func (acl *ACL) AsTableWithPIN() *lib.Table {
	header, data := acl.asTableWithPIN()

	return &lib.Table{
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

func (acl *ACL) ToTSVWithPIN(f io.Writer) error {
	header, data := acl.asTableWithPIN()

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

		header = append(header, acl.doors...)

		sort.SliceStable(acl.records, func(i, j int) bool { return acl.records[i].CardNumber < acl.records[j].CardNumber })

		for _, r := range acl.records {
			row := []string{
				fmt.Sprintf("%v", r.CardNumber),
				fmt.Sprintf("%v", r.StartDate),
				fmt.Sprintf("%v", r.EndDate),
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
							if vv {
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

				for k := range r.Revoked {
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

func (acl *ACL) asTableWithPIN() ([]string, [][]string) {
	header := []string{}
	data := [][]string{}

	if acl != nil {
		header = append(header, []string{
			"Card Number",
			"PIN",
			"From",
			"To",
		}...)

		header = append(header, acl.doors...)

		sort.SliceStable(acl.records, func(i, j int) bool { return acl.records[i].CardNumber < acl.records[j].CardNumber })

		for _, r := range acl.records {
			var pin string

			if r.PIN != 0 {
				pin = fmt.Sprintf("%v", r.PIN)
			} else {
				pin = ""
			}

			row := []string{
				fmt.Sprintf("%v", r.CardNumber),
				pin,
				fmt.Sprintf("%v", r.StartDate),
				fmt.Sprintf("%v", r.EndDate),
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
							if vv {
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

				for k := range r.Revoked {
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

func startOfYear() core.Date {
	now := time.Now()
	year := now.Year()
	month := time.January
	day := 1

	return core.ToDate(year, month, day)
}

func endOfYear() core.Date {
	now := time.Now()
	year := now.Year()
	month := time.December
	day := 31

	return core.ToDate(year, month, day)
	// return time.Date(, time.December, 31, 23, 59, 59, 0, time.Local)
}
