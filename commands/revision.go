package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Hashable interface {
	Hash() string
}

type versionInfo struct {
	AccountID uint32     `json:"account-id"`
	Timestamp *time.Time `json:"timestamp"`
	Hashes    struct {
		Members string `json:"members,omitempty"`
		Rules   string `json:"rules,omitempty"`
		ACL     string `json:"acl,omitempty"`
	} `json:"hashes"`
}

func getVersionInfo(workdir string, accountID uint32) versionInfo {
	file := filepath.Join(workdir, ".wild-apricot", fmt.Sprintf("%v.version", accountID))
	bytes, err := os.ReadFile(file)
	if err != nil {
		return versionInfo{}
	}

	v := versionInfo{}
	if err := json.Unmarshal(bytes, &v); err != nil {
		return versionInfo{}
	}

	if v.AccountID != accountID {
		return versionInfo{}
	}

	return v
}

func storeVersionInfo(workdir string, accountID uint32, timestamp time.Time, members, rules, acl Hashable) error {
	v := versionInfo{
		AccountID: accountID,
		Timestamp: &timestamp,
		Hashes: struct {
			Members string `json:"members,omitempty"`
			Rules   string `json:"rules,omitempty"`
			ACL     string `json:"acl,omitempty"`
		}{
			Members: members.Hash(),
			Rules:   rules.Hash(),
			ACL:     acl.Hash(),
		},
	}

	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	file := filepath.Join(workdir, ".wild-apricot", fmt.Sprintf("%v.version", accountID))
	bytes = append(bytes, []byte("\n")...)

	if err := os.WriteFile(file, bytes, 0644); err != nil {
		return err
	}

	return nil
}
