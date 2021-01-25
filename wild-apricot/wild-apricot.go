package wildapricot

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type permission struct {
	AccountID         uint32   `json:"AccountId"`
	SecurityProfileId uint32   `json:"SecurityProfileId"`
	AvailableScopes   []string `json:"AvailableScopes"`
}

type authorisation struct {
	AccessToken  string       `json:"access_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"`
	RefreshToken string       `json:"refresh_token"`
	Permissions  []permission `json:"Permissions"`
}

func Authorize(apiKey string) (string, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	auth := base64.StdEncoding.EncodeToString([]byte("APIKEY:" + apiKey))

	form := url.Values{
		"grant_type": []string{"client_credentials"},
		"scope":      []string{"auto"},
	}

	rq, err := http.NewRequest("POST", "https://oauth.wildapricot.org/auth/token", strings.NewReader(form.Encode()))
	rq.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rq.Header.Set("Accepts", "application/json")

	response, err := client.Do(rq)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Authorization request failed (%s)", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	authx := authorisation{}

	if err := json.Unmarshal(body, &authx); err != nil {
		return "", err
	}

	return authx.AccessToken, nil
}

func GetContacts(accountId uint32, token string) ([]Contact, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf("https://api.wildapricot.org/v2/accounts/%[1]v/contacts?$async=false", accountId)
	rq, err := http.NewRequest("GET", url, nil)

	rq.Header.Set("Accept", "application/json")
	rq.Header.Set("Authorization", "Bearer "+token)
	response, err := client.Do(rq)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	contacts := struct {
		Contacts []Contact `json:"Contacts"`
	}{}

	if err := json.Unmarshal(body, &contacts); err != nil {
		return nil, err
	}

	//	contacts := []Contact{}
	//	for _, c := range data.Contacts {
	//		contact := Contact{
	//			ID:    c.ID,
	//			Name:  fmt.Sprintf("%[1]s %[2]s", c.FirstName, c.LastName),
	//			Email: c.Email,
	//		}
	//
	//		if c.Enabled && strings.ToLower(c.Status) == "active" {
	//			contact.Active = true
	//		}
	//
	//		for _, f := range c.Fields {
	//			switch {
	//			case strings.ToLower(f.SystemCode) == "issuspendedmember":
	//				if v, ok := f.Value.(bool); ok {
	//					contact.Suspended = v
	//				}
	//
	//			case strings.ToLower(f.SystemCode) == "membersince":
	//				if v, ok := f.Value.(string); ok {
	//					if d, err := time.Parse("2006-01-02T15:04:05-07:00", v); err != nil {
	//						return nil, err
	//					} else {
	//						contact.MemberSince = &d
	//					}
	//				}
	//
	//			case strings.ToLower(f.SystemCode) == "renewaldue":
	//				if v, ok := f.Value.(string); ok {
	//					if d, err := time.Parse("2006-01-02T15:04:05", v); err != nil {
	//						return nil, err
	//					} else {
	//						contact.Renew = &d
	//					}
	//				}
	//			}
	//		}
	//
	//		contacts = append(contacts, contact)
	//	}

	return contacts.Contacts, nil
}

func GetMemberGroups(accountId uint32, token string) ([]Group, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf("https://api.wildapricot.org/v2.2/accounts/%[1]v/membergroups", accountId)
	rq, err := http.NewRequest("GET", url, nil)

	rq.Header.Set("Accept", "application/json")
	rq.Header.Set("Authorization", "Bearer "+token)
	response, err := client.Do(rq)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	data := []group{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	groups := []Group{}
	for _, g := range data {
		groups = append(groups, Group{
			ID:   g.ID,
			Name: g.Name,
		})
	}

	return groups, nil
}
