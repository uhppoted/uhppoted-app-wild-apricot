package acl

import (
	"reflect"
	"testing"
	"time"
)

func TestGrantWithDoorsOnly(t *testing.T) {
	expected := record{
		Name:       "Harry Potter",
		CardNumber: 6000001,
		StartDate:  startOfYear(),
		EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
		Granted: map[string]interface{}{
			"dungeon":    true,
			"greathall":  true,
			"gryffindor": true,
		},
		Revoked: map[string]struct{}{},
	}

	r := record{
		Name:       "Harry Potter",
		CardNumber: 6000001,
		StartDate:  startOfYear(),
		EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
		Granted:    map[string]interface{}{},
		Revoked:    map[string]struct{}{},
	}

	r.Grant("Great Hall", "Dungeon", "Gryffindor")

	if !reflect.DeepEqual(r, expected) {
		t.Errorf("'grant' failed\n   expected: %v\n   got:     %v", expected, r)
	}
}

func TestGrantWithDoorProfile(t *testing.T) {
	expected := record{
		Name:       "Harry Potter",
		CardNumber: 6000001,
		StartDate:  startOfYear(),
		EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
		Granted: map[string]interface{}{
			"dungeon":    29,
			"greathall":  true,
			"gryffindor": true,
		},
		Revoked: map[string]struct{}{},
	}

	r := record{
		Name:       "Harry Potter",
		CardNumber: 6000001,
		StartDate:  startOfYear(),
		EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
		Granted:    map[string]interface{}{},
		Revoked:    map[string]struct{}{},
	}

	r.Grant("Great Hall", "Dungeon:29", "Gryffindor")

	if !reflect.DeepEqual(r, expected) {
		t.Errorf("'grant' failed\n   expected: %v\n   got:     %v", expected, r)
	}
}

func TestGrantWithDoorAndProfile(t *testing.T) {
	expected := record{
		Name:       "Harry Potter",
		CardNumber: 6000001,
		StartDate:  startOfYear(),
		EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
		Granted: map[string]interface{}{
			"dungeon": 29,
		},
		Revoked: map[string]struct{}{},
	}

	r := record{
		Name:       "Harry Potter",
		CardNumber: 6000001,
		StartDate:  startOfYear(),
		EndDate:    time.Date(2021, time.June, 30, 0, 0, 0, 0, time.Local),
		Granted:    map[string]interface{}{},
		Revoked:    map[string]struct{}{},
	}

	r.Grant("Dungeon", 29)

	if !reflect.DeepEqual(r, expected) {
		t.Errorf("'grant' failed\n   expected: %v\n   got:     %v", expected, r)
	}
}
