package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"
)

func getTimestamp(workdir string, accountID uint32) *time.Time {
	file := filepath.Join(workdir, ".wild-apricot", fmt.Sprintf("%v.version", accountID))
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	v := struct {
		AccountID uint32 `json:"account-id"`
		Timestamp string `json:"timestamp"`
	}{}

	if err := json.Unmarshal(bytes, &v); err != nil {
		return nil
	}

	if v.AccountID != accountID {
		return nil
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", v.Timestamp)
	if err != nil {
		return nil
	}

	return &timestamp
}

func storeTimestamp(workdir string, accountID uint32, timestamp time.Time) error {
	v := struct {
		AccountID uint32 `json:"account-id"`
		Timestamp string `json:"timestamp"`
	}{
		AccountID: accountID,
		Timestamp: timestamp.Format("2006-01-02 15:04:05"),
	}

	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	file := filepath.Join(workdir, ".wild-apricot", fmt.Sprintf("%v.version", accountID))
	bytes = append(bytes, []byte("\n")...)

	if err := ioutil.WriteFile(file, bytes, 0644); err != nil {
		return err
	}

	return nil
}
