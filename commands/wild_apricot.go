package commands

import (
	"fmt"
	"strings"

	"github.com/uhppoted/uhppoted-api/config"
	"github.com/uhppoted/uhppoted-app-wild-apricot/acl"
	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
	"github.com/uhppoted/uhppoted-app-wild-apricot/wild-apricot"
)

func getMembers(conf *config.Config, credentialsFile string) (*types.Members, error) {
	cardNumberField := conf.WildApricot.Fields.CardNumber
	groupDisplayOrder := strings.Split(conf.WildApricot.DisplayOrder.Groups, ",")

	credentials, err := getCredentials(credentialsFile)
	if err != nil {
		return nil, err
	}

	token, err := wildapricot.Authorize(credentials.APIKey)
	if err != nil {
		return nil, err
	}

	contacts, err := wildapricot.GetContacts(credentials.Account, token)
	if err != nil {
		return nil, err
	}

	groups, err := wildapricot.GetMemberGroups(credentials.Account, token)
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

func getGroups(credentialsFile string, displayOrder []string) (*types.Groups, error) {
	credentials, err := getCredentials(credentialsFile)
	if err != nil {
		return nil, err
	}

	token, err := wildapricot.Authorize(credentials.APIKey)
	if err != nil {
		return nil, err
	}

	memberGroups, err := wildapricot.GetMemberGroups(credentials.Account, token)
	if err != nil {
		return nil, err
	}

	groups, err := types.MakeGroupList(memberGroups, displayOrder)
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
