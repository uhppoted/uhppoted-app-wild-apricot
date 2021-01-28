package acl

import (
	"reflect"
	"testing"
	"time"

	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
)

var C1000001 = types.CardNumber(1000001)
var C6000001 = types.CardNumber(6000001)
var C6000002 = types.CardNumber(6000002)
var C2000001 = types.CardNumber(2000001)

var dumbledore = types.Member{
	ID:         57944160,
	Name:       "Albus Dumbledore",
	CardNumber: &C1000001,
	Active:     true,
	Suspended:  false,
	Registered: dateFromString("1880-02-29"),
}

var admin = types.Member{
	ID:        57940902,
	Name:      "admin",
	Active:    false,
	Suspended: false,
}

var harry = types.Member{
	ID:         57944170,
	Name:       "Harry Potter",
	CardNumber: &C6000001,
	Active:     true,
	Suspended:  false,
	Expires:    dateFromString("2021-06-30"),
}

var hermione = types.Member{
	ID:         57944920,
	Name:       "Hermione Granger",
	CardNumber: &C6000002,
	Active:     false,
	Suspended:  false,
	Registered: dateFromString("2020-06-25"),
	Expires:    dateFromString("2021-06-30"),
}

var voldemort = types.Member{
	ID:         57944165,
	Name:       "Tom Riddle",
	CardNumber: &C2000001,
	Active:     false,
	Suspended:  true,
	Registered: dateFromString("1981-07-01"),
}

var grules = `
rule StartDate "Sets the start date to the 'registered' field" {
     when
		member.HasRegistered()
	 then
         record.SetStartDate(member.Registered);
         Retract("StartDate");
}

rule EndDate "Sets the end date to the 'expires' field" {
     when
		member.HasExpires()
	 then
         record.SetEndDate(member.Expires);
         Retract("EndDate");
}

`

func TestMakeACL(t *testing.T) {
	members := types.Members{
		Members: []types.Member{dumbledore, admin, harry, hermione, voldemort}}

	expected := ACL{
		record{
			ID:         57944160,
			Name:       "Albus Dumbledore",
			CardNumber: 1000001,
			StartDate:  time.Date(1880, time.February, 29, 0, 0, 0, 0, time.Local),
			EndDate:    endOfYear().AddDate(0, 1, 0),
		},
		record{
			ID:         57944165,
			Name:       "Tom Riddle",
			CardNumber: 2000001,
			StartDate:  time.Date(1981, time.July, 1, 0, 0, 0, 0, time.Local),
			EndDate:    endOfYear().AddDate(0, 1, 0),
		},
		record{
			ID:         57944170,
			Name:       "Harry Potter",
			CardNumber: 6000001,
			StartDate:  startOfYear(),
			EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
		},
		record{
			ID:         57944920,
			Name:       "Hermione Granger",
			CardNumber: 6000002,
			StartDate:  time.Date(2020, time.June, 25, 0, 0, 0, 0, time.Local),
			EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
		},
	}

	r, err := NewRules([]byte(grules))
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	acl, err := r.MakeACL(members)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if !reflect.DeepEqual(acl, expected) {
		if len(acl) != len(expected) {
			t.Errorf("Invalid ACL - expected %v records, got %v", len(expected), len(acl))
		} else {
			for i, _ := range expected {
				if !reflect.DeepEqual(acl[i], expected[i]) {
					t.Errorf("Invalid ACL - record %v, expected:%v, got:%v", i+1, expected[i], acl[i])
				}
			}
		}
	}
}

// TODO test with duplicate cards

func dateFromString(s string) *types.Date {
	date, err := types.DateFromString(s)
	if err != nil {
		panic(err)
	}

	return date

}
