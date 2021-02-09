package wildapricot

import (
	"compress/gzip"
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

func Authorize(apiKey string, timeout time.Duration) (string, error) {
	client := http.Client{
		Timeout: timeout,
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

func GetContacts(accountId uint32, token string, timeout time.Duration) ([]Contact, error) {
	client := http.Client{
		Timeout: timeout,
	}

	parameters := url.Values{}
	parameters.Set("$async", "false")
	parameters.Add("$filter", "'Archived' eq false AND 'Member' eq true")

	uri := fmt.Sprintf("https://api.wildapricot.org/v2/accounts/%[1]v/contacts?%[2]s", accountId, parameters.Encode())

	rq, err := http.NewRequest("GET", uri, nil)
	rq.Header.Set("Authorization", "Bearer "+token)
	rq.Header.Set("Accept", "application/json")
	rq.Header.Set("Accept-Encoding", "gzip")

	response, err := client.Do(rq)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	reader := response.Body
	if strings.ToLower(response.Header.Get("Content-Encoding")) == "gzip" {
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	contacts := struct {
		Contacts []Contact `json:"Contacts"`
	}{}

	if err := json.Unmarshal(body, &contacts); err != nil {
		return nil, err
	}

	return contacts.Contacts, nil
}

func GetMembershipLevels(accountId uint32, token string, timeout time.Duration) ([]MembershipLevel, error) {
	client := http.Client{
		Timeout: timeout,
	}

	uri := fmt.Sprintf("https://api.wildapricot.org/v2.2/accounts/%[1]v/membershiplevels", accountId)

	rq, err := http.NewRequest("GET", uri, nil)
	rq.Header.Set("Authorization", "Bearer "+token)
	rq.Header.Set("Accept", "application/json")
	rq.Header.Set("Accept-Encoding", "gzip")

	response, err := client.Do(rq)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	reader := response.Body
	if strings.ToLower(response.Header.Get("Content-Encoding")) == "gzip" {
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	fmt.Printf(">>>>>>>>>>>> %s\n", string(body))

	levels := []MembershipLevel{}
	if err := json.Unmarshal(body, &levels); err != nil {
		return nil, err
	}

	return levels, nil
}

func GetMemberGroups(accountId uint32, token string, timeout time.Duration) ([]MemberGroup, error) {
	client := http.Client{
		Timeout: timeout,
	}

	uri := fmt.Sprintf("https://api.wildapricot.org/v2.2/accounts/%[1]v/membergroups", accountId)

	rq, err := http.NewRequest("GET", uri, nil)
	rq.Header.Set("Authorization", "Bearer "+token)
	rq.Header.Set("Accept", "application/json")
	rq.Header.Set("Accept-Encoding", "gzip")

	response, err := client.Do(rq)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	reader := response.Body
	if strings.ToLower(response.Header.Get("Content-Encoding")) == "gzip" {
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return nil, err
		}
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	groups := []MemberGroup{}
	if err := json.Unmarshal(body, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

func GetUpdated(accountId uint32, token string, timestamp time.Time) (int, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	parameters := url.Values{}
	parameters.Set("$async", "false")
	parameters.Add("$filter", "'Archived' eq false AND 'Member' eq true AND 'Profile last updated' ge "+timestamp.Format("2006-01-02T15:04:05.000-07:00"))
	parameters.Add("$count", "true")

	uri := fmt.Sprintf("https://api.wildapricot.org/v2/accounts/%[1]v/contacts?%[2]s", accountId, parameters.Encode())

	rq, err := http.NewRequest("GET", uri, nil)
	rq.Header.Set("Accept", "application/json")
	rq.Header.Set("Authorization", "Bearer "+token)

	response, err := client.Do(rq)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	count := struct {
		Count int `json:"Count"`
	}{}

	if err := json.Unmarshal(body, &count); err != nil {
		return 0, err
	}

	return count.Count, nil
}
