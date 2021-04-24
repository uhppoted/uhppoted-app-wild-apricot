package commands

import (
	"reflect"
	"testing"
)

func TestCredentials(t *testing.T) {
	expected := credentials{
		AccountID: 135790,
		APIKey:    "7263hfaka9hha7d73nakd929na1nnx",
	}

	c, err := getCredentials("test_credentials.json")
	if err != nil {
		t.Fatalf("Unexpected error reading credentials (%v)", err)
	} else if c == nil {
		t.Fatalf("Invalid credentials (%v)", c)
	}

	if !reflect.DeepEqual(*c, expected) {
		t.Errorf("Incorrect credentials:\n   expected:%v,\n   got:     %v", expected, *c)
	}
}

func TestCredentialsWithMissingFile(t *testing.T) {
	_, err := getCredentials("test_missing-file.json")
	if err == nil {
		t.Errorf("Expected error reading invalid credentials, got %v", err)
	}
}

func TestCredentialsWithMissingAccountID(t *testing.T) {
	_, err := getCredentials("test_invalid-account-id.json")
	if err == nil {
		t.Errorf("Expected error reading invalid credentials, got %v", err)
	}
}

func TestCredentialsWithMissingAPIKey(t *testing.T) {
	_, err := getCredentials("test_invalid-api-key.json")
	if err == nil {
		t.Errorf("Expected error reading invalid credentials, got %v", err)
	}
}
