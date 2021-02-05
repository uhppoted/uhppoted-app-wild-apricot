package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/uhppoted/uhppoted-api/config"
	"github.com/uhppoted/uhppoted-app-wild-apricot/acl"
	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
	"github.com/uhppoted/uhppoted-app-wild-apricot/wild-apricot"
)

func revised(credentials *credentials, timestamp *time.Time) (bool, error) {
	token, err := wildapricot.Authorize(credentials.APIKey)
	if err != nil {
		return false, err
	}

	if timestamp == nil {
		return true, nil
	}

	t := timestamp.Truncate(1 * time.Second).Add(1 * time.Second)

	N, err := wildapricot.GetUpdated(credentials.AccountID, token, t)
	if err != nil {
		return false, err
	}

	return N > 0, nil
}

func getMembers(conf *config.Config, credentials *credentials) (*types.Members, error) {
	cardNumberField := conf.WildApricot.Fields.CardNumber
	groupDisplayOrder := strings.Split(conf.WildApricot.DisplayOrder.Groups, ",")

	token, err := wildapricot.Authorize(credentials.APIKey)
	if err != nil {
		return nil, err
	}

	contacts, err := wildapricot.GetContacts(credentials.AccountID, token)
	if err != nil {
		return nil, err
	}

	groups, err := wildapricot.GetMemberGroups(credentials.AccountID, token)
	if err != nil {
		return nil, err
	}

	members, err := types.MakeMemberList(contacts, groups, cardNumberField, groupDisplayOrder)
	if err != nil {
		return nil, err
	} else if members == nil {
		return nil, fmt.Errorf("Invalid members list")
	}

	return members, nil
}

func getGroups(conf *config.Config, credentials *credentials) (*types.Groups, error) {
	groupDisplayOrder := strings.Split(conf.WildApricot.DisplayOrder.Groups, ",")

	token, err := wildapricot.Authorize(credentials.APIKey)
	if err != nil {
		return nil, err
	}

	memberGroups, err := wildapricot.GetMemberGroups(credentials.AccountID, token)
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
