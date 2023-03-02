package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

type credentials struct {
	AccountID uint32 `json:"account-id"`
	APIKey    string `json:"api-key"`
}

func getCredentials(file string) (*credentials, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve credentials (%v)", err)
	}

	c := struct {
		AccountID *uint32 `json:"account-id"`
		APIKey    string  `json:"api-key"`
	}{}

	if err := json.Unmarshal(bytes, &c); err != nil {
		return nil, fmt.Errorf("unable to retrieve credentials (%v)", err)
	}

	if c.AccountID == nil {
		return nil, fmt.Errorf("invalid credentials (missing Wild Apricot account ID)")
	}

	if matched, err := regexp.MatchString("[a-zA-Z0-9]+", c.APIKey); !matched || err != nil {
		return nil, fmt.Errorf("invalid credentials (invalid Wild Apricot API key)")
	}

	return &credentials{
		AccountID: *c.AccountID,
		APIKey:    c.APIKey,
	}, nil
}
