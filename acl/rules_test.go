package acl

import (
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
)

var C1000001 = types.CardNumber(1000001)
var C6000001 = types.CardNumber(6000001)
var C6000002 = types.CardNumber(6000002)
var C2000001 = types.CardNumber(2000001)

var dumbledore = types.Member{
	Name:       "Albus Dumbledore",
	CardNumber: &C1000001,
	Active:     true,
	Suspended:  false,
	Registered: dateFromString("1880-02-29"),
	Membership: types.Membership{
		ID:   25342355,
		Name: "Staff",
	},
}

var admin = types.Member{
	Name:      "admin",
	Active:    false,
	Suspended: false,
}

var harry = types.Member{
	Name:       "Harry Potter",
	CardNumber: &C6000001,
	Active:     true,
	Suspended:  false,
	Expires:    dateFromString("2021-06-30"),
	Membership: types.Membership{
		ID:   545454,
		Name: "Stduent",
	},
}

var hermione = types.Member{
	Name:       "Hermione Granger",
	CardNumber: &C6000002,
	Active:     false,
	Suspended:  false,
	Registered: dateFromString("2020-06-25"),
	Expires:    dateFromString("2021-06-30"),
	Membership: types.Membership{
		ID:   545454,
		Name: "Stduent",
	},
}

var voldemort = types.Member{
	Name:       "Tom Riddle",
	CardNumber: &C2000001,
	Active:     false,
	Suspended:  true,
	Registered: dateFromString("1981-07-01"),
	Membership: types.Membership{
		ID:   7777777,
		Name: "Alumni",
	},
}

var grules = `
rule StartDate "Sets the start date to the 'registered' field" {
     when
		member.HasRegistered()
	 then
         permissions.SetStartDate(member.Registered);
         Retract("StartDate");
}

rule EndDate "Sets the end date to the 'expires' field" {
     when
		member.HasExpires()
	 then
         permissions.SetEndDate(member.Expires);
         Retract("EndDate");
}
`

var grant = `
rule Grant "Grants permission to the Whomping Willow" {
     when
		member.HasCardNumber(6000001)
	 then
         permissions.Grant("Whomping Willow");
         Retract("Grant");
}
`

var revoke = `
rule Revoke "Revokes permission to the Whomping Willow" {
     when
		member.HasCardNumber(6000001)
	 then
         permissions.Revoke("Whomping Willow");
         Retract("Revoke");
}
`

func TestMakeACL(t *testing.T) {
	members := types.Members{
		Members: []types.Member{dumbledore, admin, harry, hermione, voldemort},
	}

	doors := []string{}

	expected := ACL{
		records: []record{
			record{
				Name:       "Albus Dumbledore",
				CardNumber: 1000001,
				StartDate:  time.Date(1880, time.February, 29, 0, 0, 0, 0, time.Local),
				EndDate:    endOfYear().AddDate(0, 1, 0),
				Granted:    map[string]struct{}{},
				Revoked:    map[string]struct{}{},
			},
			record{
				Name:       "Tom Riddle",
				CardNumber: 2000001,
				StartDate:  time.Date(1981, time.July, 1, 0, 0, 0, 0, time.Local),
				EndDate:    endOfYear().AddDate(0, 1, 0),
				Granted:    map[string]struct{}{},
				Revoked:    map[string]struct{}{},
			},
			record{
				Name:       "Harry Potter",
				CardNumber: 6000001,
				StartDate:  startOfYear(),
				EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
				Granted:    map[string]struct{}{},
				Revoked:    map[string]struct{}{},
			},
			record{
				Name:       "Hermione Granger",
				CardNumber: 6000002,
				StartDate:  time.Date(2020, time.June, 25, 0, 0, 0, 0, time.Local),
				EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
				Granted:    map[string]struct{}{},
				Revoked:    map[string]struct{}{},
			},
		},
	}

	r, err := NewRules([]byte(grules), true)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	acl, err := r.MakeACL(members, doors)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	sort.SliceStable(expected.records, func(i, j int) bool { return expected.records[i].CardNumber < expected.records[j].CardNumber })
	sort.SliceStable(acl.records, func(i, j int) bool { return acl.records[i].CardNumber < acl.records[j].CardNumber })

	if !reflect.DeepEqual(acl, expected) {
		if len(acl.records) != len(expected.records) {
			t.Errorf("Invalid ACL - expected %v records, got %v", len(expected.records), len(acl.records))
		} else {
			for i, _ := range expected.records {
				compare(acl.records[i], expected.records[i], t)
			}
		}
	}
}

func TestMakeACLWithDuplicateCards(t *testing.T) {
	members := types.Members{
		Members: []types.Member{
			dumbledore,
			admin,
			harry,
			hermione,
			voldemort,
			types.Member{
				Name:       "Aberforth Dumbledore",
				CardNumber: &C1000001,
				Active:     true,
				Suspended:  false,
				Registered: dateFromString("2001-02-28"),
				Membership: types.Membership{
					ID:   25342355,
					Name: "Other",
				},
			},
		},
	}

	doors := []string{}

	expected := ACL{
		records: []record{
			record{
				Name:       "Albus Dumbledore",
				CardNumber: 1000001,
				StartDate:  time.Date(1880, time.February, 29, 0, 0, 0, 0, time.Local),
				EndDate:    endOfYear().AddDate(0, 1, 0),
				Granted:    map[string]struct{}{},
				Revoked:    map[string]struct{}{},
			},
			record{
				Name:       "Tom Riddle",
				CardNumber: 2000001,
				StartDate:  time.Date(1981, time.July, 1, 0, 0, 0, 0, time.Local),
				EndDate:    endOfYear().AddDate(0, 1, 0),
				Granted:    map[string]struct{}{},
				Revoked:    map[string]struct{}{},
			},
			record{
				Name:       "Harry Potter",
				CardNumber: 6000001,
				StartDate:  startOfYear(),
				EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
				Granted:    map[string]struct{}{},
				Revoked:    map[string]struct{}{},
			},
			record{
				Name:       "Hermione Granger",
				CardNumber: 6000002,
				StartDate:  time.Date(2020, time.June, 25, 0, 0, 0, 0, time.Local),
				EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
				Granted:    map[string]struct{}{},
				Revoked:    map[string]struct{}{},
			},
			record{
				Name:       "Aberforth Dumbledore",
				CardNumber: 1000001,
				StartDate:  time.Date(2001, time.February, 28, 0, 0, 0, 0, time.Local),
				EndDate:    endOfYear().AddDate(0, 1, 0),
				Granted:    map[string]struct{}{},
				Revoked:    map[string]struct{}{},
			},
		},
	}

	r, err := NewRules([]byte(grules), true)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	acl, err := r.MakeACL(members, doors)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	sort.SliceStable(expected.records, func(i, j int) bool { return expected.records[i].CardNumber < expected.records[j].CardNumber })
	sort.SliceStable(acl.records, func(i, j int) bool { return acl.records[i].CardNumber < acl.records[j].CardNumber })

	if !reflect.DeepEqual(acl, expected) {
		if len(acl.records) != len(expected.records) {
			t.Errorf("Invalid ACL - expected %v records, got %v", len(expected.records), len(acl.records))
		} else {
			for i, _ := range expected.records {
				compare(acl.records[i], expected.records[i], t)
			}
		}
	}
}

func TestGrant(t *testing.T) {
	members := types.Members{
		Members: []types.Member{harry},
	}

	doors := []string{}

	expected := ACL{
		records: []record{
			record{
				Name:       "Harry Potter",
				CardNumber: 6000001,
				StartDate:  startOfYear(),
				EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
				Granted: map[string]struct{}{
					"whompingwillow": struct{}{},
				},
				Revoked: map[string]struct{}{},
			},
		},
	}

	r, err := NewRules([]byte(grules+grant), true)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	acl, err := r.MakeACL(members, doors)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if len(acl.records) != len(expected.records) {
		t.Errorf("Invalid ACL - expected %v records, got %v", len(expected.records), len(acl.records))
	} else {
		for i, _ := range expected.records {
			compare(acl.records[i], expected.records[i], t)
		}
	}
}

func TestVariadicGrant(t *testing.T) {
	members := types.Members{
		Members: []types.Member{harry},
	}

	doors := []string{}

	expected := ACL{
		records: []record{
			record{
				Name:       "Harry Potter",
				CardNumber: 6000001,
				StartDate:  startOfYear(),
				EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
				Granted: map[string]struct{}{
					"whompingwillow": struct{}{},
					"gryffindor":     struct{}{},
					"greathall":      struct{}{},
				},
				Revoked: map[string]struct{}{},
			},
		},
	}

	grant := `
rule Grant "Grants permission to the Whomping Willow" {
     when
		member.HasCardNumber(6000001)
	 then
         permissions.Grant("Whomping Willow", "Gryffindor", "Great Hall");
         Retract("Grant");
}
`

	r, err := NewRules([]byte(grules+grant), true)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	acl, err := r.MakeACL(members, doors)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if len(acl.records) != len(expected.records) {
		t.Errorf("Invalid ACL - expected %v records, got %v", len(expected.records), len(acl.records))
	} else {
		for i, _ := range expected.records {
			compare(acl.records[i], expected.records[i], t)
		}
	}
}

func TestRevoke(t *testing.T) {
	members := types.Members{
		Members: []types.Member{harry},
	}

	doors := []string{}

	expected := ACL{
		records: []record{
			record{
				Name:       "Harry Potter",
				CardNumber: 6000001,
				StartDate:  startOfYear(),
				EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
				Granted:    map[string]struct{}{},
				Revoked: map[string]struct{}{
					"whompingwillow": struct{}{},
				},
			},
		},
	}

	r, err := NewRules([]byte(grules+revoke), true)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	acl, err := r.MakeACL(members, doors)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if len(acl.records) != len(expected.records) {
		t.Errorf("Invalid ACL - expected %v records, got %v", len(expected.records), len(acl.records))
	} else {
		for i, _ := range expected.records {
			compare(acl.records[i], expected.records[i], t)
		}
	}
}

func TestVariadicRevoke(t *testing.T) {
	members := types.Members{
		Members: []types.Member{harry},
	}

	doors := []string{}

	expected := ACL{
		records: []record{
			record{
				Name:       "Harry Potter",
				CardNumber: 6000001,
				StartDate:  startOfYear(),
				EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
				Granted:    map[string]struct{}{},
				Revoked: map[string]struct{}{
					"whompingwillow": struct{}{},
					"dungeon":        struct{}{},
					"hogsmeade":      struct{}{},
				},
			},
		},
	}

	revoke := `
rule Revoke "Revokes permission to the Whomping Willow" {
     when
		member.HasCardNumber(6000001)
	 then
         permissions.Revoke("Whomping Willow", "Dungeon", "Hogsmeade");
         Retract("Revoke");
}
`

	r, err := NewRules([]byte(grules+revoke), true)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	acl, err := r.MakeACL(members, doors)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if len(acl.records) != len(expected.records) {
		t.Errorf("Invalid ACL - expected %v records, got %v", len(expected.records), len(acl.records))
	} else {
		for i, _ := range expected.records {
			compare(acl.records[i], expected.records[i], t)
		}
	}
}

func TestGrantAndRevoke(t *testing.T) {
	members := types.Members{
		Members: []types.Member{harry},
	}

	doors := []string{}

	expected := ACL{
		records: []record{
			record{
				Name:       "Harry Potter",
				CardNumber: 6000001,
				StartDate:  startOfYear(),
				EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
				Granted: map[string]struct{}{
					"whompingwillow": struct{}{},
				},
				Revoked: map[string]struct{}{
					"whompingwillow": struct{}{},
				},
			},
		},
	}

	r, err := NewRules([]byte(grules+grant+revoke), true)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	acl, err := r.MakeACL(members, doors)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if len(acl.records) != len(expected.records) {
		t.Errorf("Invalid ACL - expected %v records, got %v", len(expected.records), len(acl.records))
	} else {
		for i, _ := range expected.records {
			compare(acl.records[i], expected.records[i], t)
		}
	}
}

func compare(r, expected record, t *testing.T) {
	if reflect.DeepEqual(r, expected) {
		return
	}

	if r.Name != expected.Name {
		t.Errorf("Invalid ACL record 'name' - expected:%v, got:%v", r.Name, expected.Name)
	}

	if r.CardNumber != expected.CardNumber {
		t.Errorf("Invalid ACL record 'card number' - expected:%v, got:%v", r.CardNumber, expected.CardNumber)
	}

	if r.StartDate.Format("2006-01-02") != expected.StartDate.Format("2006-01-02") {
		t.Errorf("Invalid ACL record 'start date' - expected:%v, got:%v", r.StartDate.Format("2006-01-02"), expected.StartDate.Format("2006-01-02"))
	}

	if r.EndDate.Format("2006-01-02") != expected.EndDate.Format("2006-01-02") {
		t.Errorf("Invalid ACL record 'end date' - expected:%v, got:%v", r.EndDate.Format("2006-01-02"), expected.EndDate.Format("2006-01-02"))
	}

	if !reflect.DeepEqual(r.Granted, expected.Granted) {
		t.Errorf("Invalid ACL record 'granted' - expected:%#v, got:%#v", r.Granted, expected.Granted)
	}

	if !reflect.DeepEqual(r.Revoked, expected.Granted) {
		t.Errorf("Invalid ACL record 'revoked' - expected:%#v, got:%#v", r.Revoked, expected.Revoked)
	}

	t.Errorf("Invalid ACL record - expected:%v, got:%v", expected, r)
}

func dateFromString(s string) *types.Date {
	date, err := types.DateFromString(s)
	if err != nil {
		panic(err)
	}

	return date

}
