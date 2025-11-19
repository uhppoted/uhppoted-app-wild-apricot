package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/uhppoted/uhppoted-app-wild-apricot/acl"
	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
	"github.com/uhppoted/uhppoted-app-wild-apricot/wild-apricot"
	"github.com/uhppoted/uhppoted-lib/config"
)

func revised(conf *config.Config, credentials *credentials, timestamp *time.Time) (bool, error) {
	token, err := wildapricot.Authorize(credentials.APIKey, conf.WildApricot.HTTP.ClientTimeout)
	if err != nil {
		return false, err
	}

	if timestamp == nil {
		return true, nil
	}

	t := timestamp.Truncate(1 * time.Second)
	timeout := conf.WildApricot.HTTP.ClientTimeout
	retries := conf.WildApricot.HTTP.Retries
	delay := conf.WildApricot.HTTP.RetryDelay

	N, err := wildapricot.GetUpdated(credentials.AccountID, token, t, timeout, retries, delay)
	if err != nil {
		return false, err
	}

	infof("Updated records: %v", N)

	return N > 0, nil
}

func getMembers(conf *config.Config, credentials *credentials) (*types.Members, error) {
	timeout := conf.WildApricot.HTTP.ClientTimeout

	cardNumberField := conf.WildApricot.Fields.CardNumber
	pinField := conf.WildApricot.Fields.PIN
	facilityCode := conf.WildApricot.FacilityCode
	groupDisplayOrder := strings.Split(conf.WildApricot.DisplayOrder.Groups, ",")

	token, err := wildapricot.Authorize(credentials.APIKey, timeout)
	if err != nil {
		return nil, err
	}

	api := wildapricot.API{
		PageSize: conf.WildApricot.HTTP.PageSize,
		MaxPages: conf.WildApricot.HTTP.MaxPages,
		Timeout:  conf.WildApricot.HTTP.ClientTimeout,
		Retries:  conf.WildApricot.HTTP.Retries,
		Delay:    conf.WildApricot.HTTP.RetryDelay,
	}

	contacts, err := wildapricot.GetContacts(credentials.AccountID, token, api)
	if err != nil {
		return nil, err
	}

	groups, err := wildapricot.GetMemberGroups(credentials.AccountID, token, api)
	if err != nil {
		return nil, err
	}

	members, errors := types.MakeMemberList(contacts, groups, cardNumberField, pinField, facilityCode, groupDisplayOrder)
	for _, err := range errors {
		warnf("%v", err.Error())
	}

	if members == nil {
		return nil, fmt.Errorf("invalid members list")
	}

	return members, nil
}

func getGroups(conf *config.Config, credentials *credentials) (*types.Groups, error) {
	timeout := conf.WildApricot.HTTP.ClientTimeout

	groupDisplayOrder := strings.Split(conf.WildApricot.DisplayOrder.Groups, ",")

	token, err := wildapricot.Authorize(credentials.APIKey, timeout)
	if err != nil {
		return nil, err
	}

	api := wildapricot.API{
		PageSize: conf.WildApricot.HTTP.PageSize,
		MaxPages: conf.WildApricot.HTTP.MaxPages,
		Timeout:  conf.WildApricot.HTTP.ClientTimeout,
		Retries:  conf.WildApricot.HTTP.Retries,
		Delay:    conf.WildApricot.HTTP.RetryDelay,
	}

	memberGroups, err := wildapricot.GetMemberGroups(credentials.AccountID, token, api)
	if err != nil {
		return nil, err
	}

	groups, err := types.MakeGroupList(memberGroups, groupDisplayOrder)
	if err != nil {
		return nil, err
	} else if groups == nil {
		return nil, fmt.Errorf("invalid groups list")
	}

	return groups, nil
}

// Ref. https://github.com/uhppoted/uhppoted-app-wild-apricot/issues/2
func getRules(uri string, workdir string, dbg bool) (*acl.Rules, error) {
	ruleset, err := fetch(uri)
	if err != nil {
		return nil, err
	}

	if dbg {
		filename := time.Now().Format("RULES 2006-01-02 15:04:05.grl")
		path := filepath.Join(os.TempDir(), filename)
		if f, err := os.Create(path); err != nil {
			warnf("%v", err)
		} else {
			f.Write(ruleset)
			f.Close()
			debugf("Stashed rules list in file %s", path)
		}
	}

	if rules, err := acl.NewRules(ruleset, dbg); err != nil {
		warnf("%v", err)
	} else {
		stash(ruleset, filepath.Join(workdir, "wild-apricot.grl"))

		return rules, nil
	}

	// ... try load cached rules
	stashed := filepath.Join(workdir, "wild-apricot.grl")
	if ruleset, err := os.ReadFile(stashed); err != nil {
		return nil, err
	} else {
		warnf("Using stashed 'grules' file (%v)", stashed)
		return acl.NewRules(ruleset, dbg)
	}
}

func stash(bytes []byte, file string) {
	if f, err := os.Create(file); err != nil {
		warnf("Error creating stashed 'grules' file (%v)", err)
	} else {
		f.Write(bytes)
		f.Close()
	}

	infof("Stashed downloaded 'grules' file to %v", file)
}
