package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func getVersion(workdir string, accountID uint32) string {
	file := filepath.Join(workdir, ".wild-apricot", fmt.Sprintf("%v.version", accountID))
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}

	v := struct {
		AccountID uint32 `json:"account-id"`
		Version   string `json:"version"`
	}{}

	if err := json.Unmarshal(bytes, &v); err != nil {
		return ""
	}

	if v.AccountID != accountID {
		return ""
	}

	return clean(v.Version)
}
