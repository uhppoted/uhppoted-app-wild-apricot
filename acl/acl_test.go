package acl

import (
	"reflect"
	"testing"
	"time"

	api "github.com/uhppoted/uhppoted-api/acl"
)

func TestAsTable(t *testing.T) {
	acl := ACL{
		doors: []string{
			"Great Hall",
			"Whomping Willow",
			"Dungeon",
			"Hogsmeade",
		},

		records: []record{
			record{
				Name:       "Albus Dumbledore",
				CardNumber: 1000001,
				StartDate:  time.Date(1880, time.February, 29, 0, 0, 0, 0, time.Local),
				EndDate:    endOfYear().AddDate(0, 1, 0),
				Granted: map[string]interface{}{
					"Great Hall":      true,
					"Whomping Willow": true,
					"Dungeon":         29,
					"Hogsmeade":       true,
				},
				Revoked: map[string]struct{}{},
			},
			record{
				Name:       "Tom Riddle",
				CardNumber: 2000001,
				StartDate:  time.Date(1981, time.July, 1, 0, 0, 0, 0, time.Local),
				EndDate:    endOfYear().AddDate(0, 1, 0),
				Granted:    map[string]interface{}{},
				Revoked:    map[string]struct{}{},
			},
			record{
				Name:       "Harry Potter",
				CardNumber: 6000001,
				StartDate:  startOfYear(),
				EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
				Granted: map[string]interface{}{
					"Great Hall": true,
					"Hogsmeade":  true,
				},
				Revoked: map[string]struct{}{
					"Hogsmeade": struct{}{},
				},
			},
			record{
				Name:       "Hermione Granger",
				CardNumber: 6000002,
				StartDate:  time.Date(2020, time.June, 25, 0, 0, 0, 0, time.Local),
				EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
				Granted: map[string]interface{}{
					"Great Hall": true,
					"Hogsmeade":  true,
				},
				Revoked: map[string]struct{}{},
			},
		},
	}

	expected := api.Table{
		Header: []string{
			"Card Number",
			"From",
			"To",
			"Great Hall",
			"Whomping Willow",
			"Dungeon",
			"Hogsmeade",
		},

		Records: [][]string{
			[]string{"1000001", "1880-02-29", "2022-01-31", "Y", "Y", "29", "Y"},
			[]string{"2000001", "1981-07-01", "2022-01-31", "N", "N", "N", "N"},
			[]string{"6000001", "2021-01-01", "2021-06-30", "Y", "N", "N", "N"},
			[]string{"6000002", "2020-06-25", "2021-06-30", "Y", "N", "N", "Y"},
		},
	}

	table := acl.AsTable()

	if !reflect.DeepEqual(table.Header, expected.Header) {
		t.Errorf("Invalid ACL table header - expected:%v, got:%v", expected.Header, table.Header)
	}

	if !reflect.DeepEqual(table.Records, expected.Records) {
		for i := range expected.Records {
			if !reflect.DeepEqual(table.Records[i], expected.Records[i]) {
				t.Errorf("Invalid ACL table row\n   expected:%v\n   got:     %v", expected.Records[i], table.Records[i])
			}
		}
	}
}
