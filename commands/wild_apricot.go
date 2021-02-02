package commands

import (
	"fmt"

	"github.com/uhppoted/uhppoted-app-wild-apricot/acl"
	"github.com/uhppoted/uhppoted-app-wild-apricot/types"
	"github.com/uhppoted/uhppoted-app-wild-apricot/wild-apricot"
)

func getMembers(file string) (*types.Members, error) {
	credentials, err := getCredentials(file)
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

	members, err := types.MakeMemberList(contacts, groups)
	if err != nil {
		return nil, err
	} else if members == nil {
		return nil, fmt.Errorf("Invalid members list")
	}

	return members, nil
}

func getRules(uri string, debug bool) (*acl.Rules, error) {
	ruleset, err := fetch(uri)
	if err != nil {
		return nil, err
	}

	return acl.NewRules(ruleset, debug)
}
