package acl

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	core "github.com/uhppoted/uhppote-core/types"
	lib "github.com/uhppoted/uhppoted-lib/acl"
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
				StartDate:  core.MustParseDate("1880-02-29"),
				EndDate:    plusOneDay(endOfYear()),
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
				StartDate:  core.MustParseDate("1981-07-01"),
				EndDate:    plusOneDay(endOfYear()),
				Granted:    map[string]interface{}{},
				Revoked:    map[string]struct{}{},
			},
			record{
				Name:       "Harry Potter",
				CardNumber: 6000001,
				StartDate:  startOfYear(),
				EndDate:    core.MustParseDate("2021-06-30"),
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
				StartDate:  core.MustParseDate("2020-06-25"),
				EndDate:    core.MustParseDate("2021-06-30"),
				Granted: map[string]interface{}{
					"Great Hall": true,
					"Hogsmeade":  true,
				},
				Revoked: map[string]struct{}{},
			},
		},
	}

	year := time.Now().Year()

	expected := lib.Table{
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
			[]string{"1000001", "1880-02-29", fmt.Sprintf("%04d-01-31", year+1), "Y", "Y", "29", "Y"},
			[]string{"2000001", "1981-07-01", fmt.Sprintf("%04d-01-31", year+1), "N", "N", "N", "N"},
			[]string{"6000001", fmt.Sprintf("%04d-01-01", year), "2021-06-30", "Y", "N", "N", "N"},
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

func TestHash(t *testing.T) {
	dumbledore := record{
		Name:       "Albus Dumbledore",
		CardNumber: 1000001,
		StartDate:  core.ToDate(1880, time.February, 29),
		EndDate:    core.ToDate(2021, time.December, 31), // FIXME EndDate:    time.Date(2021, time.December, 31, 23, 59, 59, 0, time.Local).AddDate(0, 1, 0),
		Granted: map[string]interface{}{
			"Great Hall":      true,
			"Whomping Willow": true,
			"Dungeon":         29,
			"Hogsmeade":       true,
		},
		Revoked: map[string]struct{}{},
	}

	riddle := record{
		Name:       "Tom Riddle",
		CardNumber: 2000001,
		StartDate:  core.ToDate(1981, time.July, 1),
		EndDate:    core.ToDate(2021, time.December, 31), // FIXME EndDate:    time.Date(2021, time.December, 31, 23, 59, 59, 0, time.Local).AddDate(0, 1, 0),
		Granted:    map[string]interface{}{},
		Revoked:    map[string]struct{}{},
	}

	potter := []record{
		record{
			Name:       "Harry Potter",
			CardNumber: 6000001,
			StartDate:  core.ToDate(2021, time.January, 1),
			EndDate:    core.ToDate(2021, time.June, 30),
			Granted: map[string]interface{}{
				"Great Hall": true,
				"Hogsmeade":  true,
			},
			Revoked: map[string]struct{}{
				"Hogsmeade": struct{}{},
			},
		},
		record{
			Name:       "Harry Potter",
			CardNumber: 6000001,
			StartDate:  core.ToDate(2021, time.January, 1),
			EndDate:    core.ToDate(2021, time.June, 30),
			Granted: map[string]interface{}{
				"Great Hall": 29,
				"Hogsmeade":  true,
			},
			Revoked: map[string]struct{}{
				"Hogsmeade": struct{}{},
			},
		},
		record{
			Name:       "Harry Potter",
			CardNumber: 6000001,
			StartDate:  core.ToDate(2021, time.January, 1),
			EndDate:    core.ToDate(2021, time.June, 30),
			Granted: map[string]interface{}{
				"Great Hall": true,
				"Hogsmeade":  true,
			},
			Revoked: map[string]struct{}{
				"Hogsmeade":  struct{}{},
				"Great Hall": struct{}{},
			},
		},
		record{
			Name:       "Harry Potter",
			CardNumber: 6000001,
			StartDate:  core.ToDate(2021, time.January, 1),
			EndDate:    core.ToDate(2021, time.June, 30),
			Granted: map[string]interface{}{
				"Great Hall": true,
				"Hogsmeade":  true,
			},
			Revoked: map[string]struct{}{
				"Hogsmeade": struct{}{},
				"Dungeon":   struct{}{},
			},
		},
	}

	acl := ACL{
		doors:   []string{"Great Hall", "Whomping Willow", "Dungeon", "Hogsmeade"},
		records: []record{dumbledore, riddle, potter[0]},
	}

	// FIXME double check (end date has changed)
	// expected := "2257f356a9efe68827e7324d4cd68f73b0e9127a1dac65af5659cb87578ef5dc"
	expected := "d2c28d1f2559a539ce99b8f33115f374bc25e3fce483e9e8ba84ca5965741aed"
	hash := acl.Hash()
	if hash != expected {
		t.Errorf("Invalid ACL hash - expected:%v, got:%v", expected, hash)
	}

	acl = ACL{
		doors:   []string{"Great Hall", "Whompping Willow", "Dungeon", "Hogsmeade"},
		records: []record{dumbledore, riddle, potter[1]},
	}

	hash = acl.Hash()
	if hash == expected {
		t.Errorf("Invalid ACL hash - expected:%v, got:%v", "71b3ab210f0e0112691687e92ea12a49241f9dc7c2729f52916fe2060b0ec74e", hash)
	}

	acl = ACL{
		doors:   []string{"Great Hall", "Whomping Willow", "Dungeon", "Hogsmeade"},
		records: []record{dumbledore, riddle, potter[2]},
	}

	hash = acl.Hash()
	if hash == expected {
		t.Errorf("Invalid ACL hash - expected:%v, got:%v", "b829378a79576c1bb41c78d3fe23a32e310d54fa6d8d86945909440319b5480b", hash)
	}

	acl = ACL{
		doors:   []string{"Great Hall", "Whomping Willow", "Dungeon", "Hogsmeade"},
		records: []record{dumbledore, riddle, potter[3]},
	}

	hash = acl.Hash()
	if hash != expected {
		t.Errorf("Invalid ACL hash - expected:%v, got:%v", expected, hash)
	}

}
