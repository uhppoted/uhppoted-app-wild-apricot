package types

import (
	"testing"
)

func TestNormalisation(t *testing.T) {
	tests := []struct {
		s        string
		expected string
	}{
		{
			`ABCDEFGHIJKLMNOPQRSTUVWXabcdefghijklmnopqrstuvwx0123456789!@#$%^&*()-_+={}[]|\;:'"<,>.?/`,
			`abcdefghijklmnopqrstuvwxabcdefghijklmnopqrstuvwx0123456789`,
		},
	}

	for _, test := range tests {
		v := normalise(test.s)

		if v != test.expected {
			t.Errorf("incorrectly normalised '%v'\n   expected:'%v'\n   got:     '%v'", test.s, test.expected, v)
		}
	}
}
