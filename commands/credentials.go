package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type credentials struct {
	Account uint32 `json:"account"`
	APIKey  string `json:"api-key"`
}

func getCredentials(file string) (*credentials, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve credentials (%v)", err)
	}

	c := credentials{}

	if err := json.Unmarshal(bytes, &c); err != nil {
		return nil, fmt.Errorf("Unable to retrieve credentials (%v)", err)
	}

	return &c, nil
}
