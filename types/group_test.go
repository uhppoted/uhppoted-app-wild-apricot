package types

import (
	"math"
	"reflect"
	"testing"

	"github.com/uhppoted/uhppoted-app-wild-apricot/wild-apricot"
)

func TestGroupConstructor(t *testing.T) {
	mg := wildapricot.MemberGroup{
		ID:          654321,
		Name:        "Gryffindor",
		Description: "Group for Gryffindoor students",
		URL:         "https://api.wildapricot.org/v2.2/accounts/12345/MemberGroups/654321",
		Contacts:    5,
	}

	expected := Group{
		ID:    654321,
		Name:  "Gryffindor",
		index: math.MaxUint32,
	}

	group, err := NewGroup(mg)
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if !reflect.DeepEqual(group, expected) {
		t.Errorf("Invalid group - expected:%+v, got:%+v", expected, group)
	}
}
