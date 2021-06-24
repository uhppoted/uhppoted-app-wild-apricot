package commands

import (
	"fmt"
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

	info(fmt.Sprintf("Updated records: %v", N))

	return N > 0, nil
}

func getMembers(conf *config.Config, credentials *credentials) (*types.Members, error) {
	timeout := conf.WildApricot.HTTP.ClientTimeout
	retries := conf.WildApricot.HTTP.Retries
	delay := conf.WildApricot.HTTP.RetryDelay

	cardNumberField := conf.WildApricot.Fields.CardNumber
	facilityCode := conf.WildApricot.FacilityCode
	groupDisplayOrder := strings.Split(conf.WildApricot.DisplayOrder.Groups, ",")

	token, err := wildapricot.Authorize(credentials.APIKey, timeout)
	if err != nil {
		return nil, err
	}

	contacts, err := wildapricot.GetContacts(credentials.AccountID, token, timeout, retries, delay)
	if err != nil {
		return nil, err
	}

	groups, err := wildapricot.GetMemberGroups(credentials.AccountID, token, timeout, retries, delay)
	if err != nil {
		return nil, err
	}

	members, errors := types.MakeMemberList(contacts, groups, cardNumberField, facilityCode, groupDisplayOrder)
	for _, err := range errors {
		warn(err.Error())
	}

	if members == nil {
		return nil, fmt.Errorf("Invalid members list")
	}

	return members, nil
}

func getGroups(conf *config.Config, credentials *credentials) (*types.Groups, error) {
	timeout := conf.WildApricot.HTTP.ClientTimeout
	retries := conf.WildApricot.HTTP.Retries
	delay := conf.WildApricot.HTTP.RetryDelay

	groupDisplayOrder := strings.Split(conf.WildApricot.DisplayOrder.Groups, ",")

	token, err := wildapricot.Authorize(credentials.APIKey, timeout)
	if err != nil {
		return nil, err
	}

	memberGroups, err := wildapricot.GetMemberGroups(credentials.AccountID, token, timeout, retries, delay)
	if err != nil {
		return nil, err
	}

	groups, err := types.MakeGroupList(memberGroups, groupDisplayOrder)
	if err != nil {
		return nil, err
	} else if groups == nil {
		return nil, fmt.Errorf("Invalid groups list")
	}

	return groups, nil
}

func getRules(uri string, debug bool) (*acl.Rules, error) {
	ruleset, err := fetch(uri)
	if err != nil {
		return nil, err
	}

	return acl.NewRules(ruleset, debug)
}
