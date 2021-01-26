package wildapricot

import (
	"reflect"
	"testing"
)

func TestMemberGroupFlatten(t *testing.T) {
	mg := MemberGroup{
		ID:          654321,
		Name:        "Gryffindor",
		Description: "Group for Gryffindoor students",
		URL:         "https://api.wildapricot.org/v2.2/accounts/12345/MemberGroups/654321",
		Contacts:    5,
	}

	expected := map[string]interface{}{
		"id":          uint32(654321),
		"name":        "Gryffindor",
		"description": "Group for Gryffindoor students",
		"contacts":    5,
		"url":         "https://api.wildapricot.org/v2.2/accounts/12345/MemberGroups/654321",
	}

	m, err := mg.Flatten()
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Invalid flatten - expected:%#v, got:%#v", expected, m)
	}
}

func TestMemberGroupFlattenWithNil(t *testing.T) {
	var mg *MemberGroup

	expected := map[string]interface{}{}

	m, err := mg.Flatten()
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Invalid flatten - expected:%+v, got:%+v", expected, m)
	}
}
