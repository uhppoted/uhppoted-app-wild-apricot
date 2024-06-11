package acl

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"

	core "github.com/uhppoted/uhppote-core/types"

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
	Registered: core.MustParseDate("1880-02-29"),
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
	Expires:    core.MustParseDate("2021-06-30"),
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
	Registered: core.MustParseDate("2020-06-25"),
	Expires:    core.MustParseDate("2021-06-30"),
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
	Registered: core.MustParseDate("1981-07-01"),
	Membership: types.Membership{
		ID:   7777777,
		Name: "Alumni",
	},
}

var grules = `
// *** GRULES ***
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
// *** END GRULES ***
`

var grant = `
// *** GRULES ***
rule Grant "Grants permission to the Whomping Willow" {
     when
		member.HasCardNumber(6000001)
	 then
         permissions.Grant("Whomping Willow");
         Retract("Grant");
}
// *** END GRULES ***
`

var grantWithTimeProfile = `
// *** GRULES ***
rule Grant "Grants permission to the Whomping Willow" {
     when
		member.HasCardNumber(6000001)
	 then
         permissions.Grant("Whomping Willow:29");
         Retract("Grant");
}
// *** END GRULES ***
`

var grantAndTimeProfile = `
// *** GRULES ***
rule Grant "Grants permission to the Whomping Willow" {
     when
		member.HasCardNumber(6000001)
	 then
         permissions.Grant("Whomping Willow",55);
         Retract("Grant");
}
// *** END GRULES ***
`

var revoke = `
// *** GRULES ***
rule Revoke "Revokes permission to the Whomping Willow" {
     when
		member.HasCardNumber(6000001)
	 then
         permissions.Revoke("Whomping Willow");
         Retract("Revoke");
}
// *** END GRULES ***
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
				StartDate:  core.ToDate(1880, time.February, 29),
				EndDate:    plusOneDay(endOfYear()),
				Granted:    map[string]interface{}{},
				Revoked:    map[string]struct{}{},
			},
			record{
				Name:       "Tom Riddle",
				CardNumber: 2000001,
				StartDate:  core.ToDate(1981, time.July, 1),
				EndDate:    plusOneDay(endOfYear()),
				Granted:    map[string]interface{}{},
				Revoked:    map[string]struct{}{},
			},
			record{
				Name:       "Harry Potter",
				CardNumber: 6000001,
				StartDate:  startOfYear(),
				EndDate:    core.ToDate(2021, time.June, 30),
				Granted:    map[string]interface{}{},
				Revoked:    map[string]struct{}{},
			},
			record{
				Name:       "Hermione Granger",
				CardNumber: 6000002,
				StartDate:  core.ToDate(2020, time.June, 25),
				EndDate:    core.ToDate(2021, time.June, 30),
				Granted:    map[string]interface{}{},
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
			for i := range expected.records {
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
				Registered: core.MustParseDate("2001-02-28"),
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
				StartDate:  core.ToDate(1880, time.February, 29),
				EndDate:    plusOneDay(endOfYear()),
				Granted:    map[string]interface{}{},
				Revoked:    map[string]struct{}{},
			},
			record{
				Name:       "Tom Riddle",
				CardNumber: 2000001,
				StartDate:  core.ToDate(1981, time.July, 1),
				EndDate:    plusOneDay(endOfYear()),
				Granted:    map[string]interface{}{},
				Revoked:    map[string]struct{}{},
			},
			record{
				Name:       "Harry Potter",
				CardNumber: 6000001,
				StartDate:  startOfYear(),
				EndDate:    core.ToDate(2021, time.June, 30),
				Granted:    map[string]interface{}{},
				Revoked:    map[string]struct{}{},
			},
			record{
				Name:       "Hermione Granger",
				CardNumber: 6000002,
				StartDate:  core.ToDate(2020, time.June, 25),
				EndDate:    core.ToDate(2021, time.June, 30),
				Granted:    map[string]interface{}{},
				Revoked:    map[string]struct{}{},
			},
			record{
				Name:       "Aberforth Dumbledore",
				CardNumber: 1000001,
				StartDate:  core.ToDate(2001, time.February, 28),
				EndDate:    plusOneDay(endOfYear()),
				Granted:    map[string]interface{}{},
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
			for i := range expected.records {
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
				EndDate:    core.ToDate(2021, time.June, 30),
				Granted: map[string]interface{}{
					"whompingwillow": true,
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
		for i := range expected.records {
			compare(acl.records[i], expected.records[i], t)
		}
	}
}

func TestGrantWithDoorWithTimeProfile(t *testing.T) {
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
				EndDate:    core.ToDate(2021, time.June, 30),
				Granted: map[string]interface{}{
					"whompingwillow": 29,
				},
				Revoked: map[string]struct{}{},
			},
		},
	}

	r, err := NewRules([]byte(grules+grantWithTimeProfile), true)
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
		for i := range expected.records {
			compare(acl.records[i], expected.records[i], t)
		}
	}
}

func TestGrantWithDoorAndTimeProfile(t *testing.T) {
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
				EndDate:    core.ToDate(2021, time.June, 30),
				Granted: map[string]interface{}{
					"whompingwillow": 55,
				},
				Revoked: map[string]struct{}{},
			},
		},
	}

	r, err := NewRules([]byte(grules+grantAndTimeProfile), true)
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
		for i := range expected.records {
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
				EndDate:    core.ToDate(2021, time.June, 30),
				Granted: map[string]interface{}{
					"whompingwillow": true,
					"gryffindor":     true,
					"greathall":      true,
				},
				Revoked: map[string]struct{}{},
			},
		},
	}

	grant := `
// *** GRULES ***
rule Grant "Grants permission to the Whomping Willow" {
     when
		member.HasCardNumber(6000001)
	 then
         permissions.Grant("Whomping Willow", "Gryffindor", "Great Hall");
         Retract("Grant");
}
// *** END GRULES ***
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
		for i := range expected.records {
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
				EndDate:    core.ToDate(2021, time.June, 30),
				Granted:    map[string]interface{}{},
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
		for i := range expected.records {
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
				EndDate:    core.ToDate(2021, time.June, 30),
				Granted:    map[string]interface{}{},
				Revoked: map[string]struct{}{
					"whompingwillow": struct{}{},
					"dungeon":        struct{}{},
					"hogsmeade":      struct{}{},
				},
			},
		},
	}

	revoke := `
// *** GRULES ***
rule Revoke "Revokes permission to the Whomping Willow" {
     when
		member.HasCardNumber(6000001)
	 then
         permissions.Revoke("Whomping Willow", "Dungeon", "Hogsmeade");
         Retract("Revoke");
}
// *** END GRULES ***
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
		for i := range expected.records {
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
				EndDate:    core.ToDate(2021, time.June, 30),
				Granted: map[string]interface{}{
					"whompingwillow": true,
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
		for i := range expected.records {
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

	// FIXME use date.Equal
	if fmt.Sprintf("%v", r.StartDate) != fmt.Sprintf("%v", expected.StartDate) {
		t.Errorf("Invalid ACL record 'start date' - expected:%v, got:%v", r.StartDate, expected.StartDate)
	}

	// FIXME use date.Equal
	if fmt.Sprintf("%v", r.EndDate) != fmt.Sprintf("%v", expected.EndDate) {
		t.Errorf("Invalid ACL record 'end date' - expected:%v, got:%v", r.EndDate, expected.EndDate)
	}

	if !reflect.DeepEqual(r.Granted, expected.Granted) {
		t.Errorf("Invalid ACL record 'granted' - expected:%#v, got:%#v", r.Granted, expected.Granted)
	}

	if !reflect.DeepEqual(r.Revoked, expected.Granted) {
		t.Errorf("Invalid ACL record 'revoked' - expected:%#v, got:%#v", r.Revoked, expected.Revoked)
	}

	t.Errorf("Invalid ACL record - expected:%v, got:%v", expected, r)
}
