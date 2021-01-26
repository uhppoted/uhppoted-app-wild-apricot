package rules

import (
	"reflect"
	"testing"

	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
)

var C1000001 = types.CardNumber(1000001)
var C6000001 = types.CardNumber(6000001)
var C6000002 = types.CardNumber(6000002)
var C2000001 = types.CardNumber(2000001)

func TestMakeACL(t *testing.T) {
	members := types.Members{
		Members: []types.Member{
			types.Member{
				ID:         57944160,
				Name:       "Albus Dumbledore",
				CardNumber: &C1000001,
				Active:     true,
				Suspended:  false},

			types.Member{
				ID:        57940902,
				Name:      "development uhppoted",
				Active:    false,
				Suspended: false},

			types.Member{
				ID:         57944170,
				Name:       "Harry Potter",
				CardNumber: &C6000001,
				Active:     true,
				Suspended:  false},

			types.Member{
				ID:         57944920,
				Name:       "Hermione Granger",
				CardNumber: &C6000002,
				Active:     false,
				Suspended:  false},

			types.Member{
				ID:         57944165,
				Name:       "Tom Riddle",
				CardNumber: &C2000001,
				Active:     false,
				Suspended:  true},
		},
	}

	expected := ACL{
		record{
			ID:         57944160,
			Name:       "Albus Dumbledore",
			CardNumber: 1000001,
		},
		record{
			ID:         57944165,
			Name:       "Tom Riddle",
			CardNumber: 2000001,
		},
		record{
			ID:         57944170,
			Name:       "Harry Potter",
			CardNumber: 6000001,
		},
		record{
			ID:         57944920,
			Name:       "Hermione Granger",
			CardNumber: 6000002,
		},
	}

	acl, err := MakeACL(members)
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
