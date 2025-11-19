package wildapricot

import (
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type API struct {
	PageSize int
	MaxPages int

	Timeout time.Duration
	Retries int
	Delay   time.Duration
}

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

const MinPageSize = 2
const MaxPageSize = 100
const MinPages = 10
const MaxPages = 50

func Authorize(apiKey string, timeout time.Duration) (string, error) {
	client := http.Client{
		Timeout: timeout,
	}

	auth := base64.StdEncoding.EncodeToString([]byte("APIKEY:" + apiKey))

	form := url.Values{
		"grant_type": []string{"client_credentials"},
		"scope":      []string{"auto"},
	}

	rq, _ := http.NewRequest("POST", "https://oauth.wildapricot.org/auth/token", strings.NewReader(form.Encode()))
	rq.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rq.Header.Set("Accepts", "application/json")

	response, err := client.Do(rq)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("authorization request failed (%s)", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	authx := authorisation{}

	if err := json.Unmarshal(body, &authx); err != nil {
		return "", err
	}

	return authx.AccessToken, nil
}

func GetContacts(accountId uint32, token string, api API) ([]Contact, error) {
	list := []Contact{}
	pageSize := uint32(api.PageSize)
	maxPages := uint32(api.MaxPages)
	pages := uint32(0)
	page := 0

	if pageSize < MinPageSize {
		pageSize = MinPageSize
	} else if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	if maxPages < MinPages {
		maxPages = MinPages
	} else if maxPages > MaxPages {
		maxPages = MaxPages
	}

	for pages < maxPages {
		if contacts, err := getContacts(accountId, token, pageSize, uint32(page), api); err != nil {
			return nil, err
		} else if len(contacts) == 0 {
			return list, nil
		} else {
			list = append(list, contacts...)
			page += len(contacts)

			if len(contacts) == 1 {
				info("retrieved %v member (total: %v)", len(contacts), len(list))
			} else {
				info("retrieved %v members (total: %v)", len(contacts), len(list))
			}
		}

		pages++
	}

	return nil, fmt.Errorf("failed to retrieve entire contact list in %v page requests", pages)
}

func getContacts(accountId uint32, token string, pageSize uint32, page uint32, api API) ([]Contact, error) {
	parameters := url.Values{}
	parameters.Set("$async", "false")
	parameters.Add("$top", fmt.Sprintf("%v", pageSize))
	parameters.Add("$skip", fmt.Sprintf("%v", page))
	parameters.Add("$filter", "'Archived' eq false AND 'Member' eq true")

	uri := fmt.Sprintf("https://api.wildapricot.org/v2/accounts/%[1]v/contacts?%[2]s", accountId, parameters.Encode())

	rq, _ := http.NewRequest("GET", uri, nil)
	rq.Header.Set("Authorization", "Bearer "+token)
	rq.Header.Set("Accept", "application/json")
	rq.Header.Set("Accept-Encoding", "gzip")

	response, err := get(rq, api.Timeout, api.Retries, api.Delay)
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

	body, err := io.ReadAll(reader)
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

func GetMemberGroups(accountId uint32, token string, api API) ([]MemberGroup, error) {
	list := []MemberGroup{}
	pageSize := uint32(api.PageSize)
	maxPages := uint32(api.MaxPages)
	pages := uint32(0)
	page := 0

	if pageSize < MinPageSize {
		pageSize = MinPageSize
	} else if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	if maxPages < MinPages {
		maxPages = MinPages
	} else if maxPages > MaxPages {
		maxPages = MaxPages
	}

	for pages < maxPages {
		if groups, err := getMemberGroups(accountId, token, pageSize, uint32(page), api); err != nil {
			return nil, err
		} else if len(groups) == 0 {
			return list, nil
		} else {
			list = append(list, groups...)
			page += len(groups)

			if len(groups) == 1 {
				info("retrieved %v group (total: %v)", len(groups), len(list))
			} else {
				info("retrieved %v groups (total: %v)", len(groups), len(list))
			}
		}

		pages++
	}

	return nil, fmt.Errorf("failed to retrieve entire group list in %v page requests", pages)

	// return getMemberGroups(accountId, token, timeout, retries, delay, 3, 0)
}

func getMemberGroups(accountId uint32, token string, pageSize uint32, page uint32, api API) ([]MemberGroup, error) {
	parameters := url.Values{}
	parameters.Set("$async", "false")
	parameters.Add("$top", fmt.Sprintf("%v", pageSize))
	parameters.Add("$skip", fmt.Sprintf("%v", page))

	uri := fmt.Sprintf("https://api.wildapricot.org/v2.2/accounts/%[1]v/membergroups?%[2]s", accountId, parameters.Encode())

	rq, _ := http.NewRequest("GET", uri, nil)
	rq.Header.Set("Authorization", "Bearer "+token)
	rq.Header.Set("Accept", "application/json")
	rq.Header.Set("Accept-Encoding", "gzip")

	response, err := get(rq, api.Timeout, api.Retries, api.Delay)
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

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	groups := []MemberGroup{}
	if err := json.Unmarshal(body, &groups); err != nil {
		return nil, err
	}

	return groups, nil
}

func GetUpdated(accountId uint32, token string, timestamp time.Time, api API) (int, error) {
	parameters := url.Values{}
	parameters.Set("$async", "false")
	parameters.Add("$filter", "'Archived' eq false AND 'Profile last updated' ge "+timestamp.Format("2006-01-02T15:04:05.000-07:00"))
	parameters.Add("$count", "true")

	uri := fmt.Sprintf("https://api.wildapricot.org/v2/accounts/%[1]v/contacts?%[2]s", accountId, parameters.Encode())

	rq, _ := http.NewRequest("GET", uri, nil)
	rq.Header.Set("Accept", "application/json")
	rq.Header.Set("Authorization", "Bearer "+token)

	response, err := get(rq, api.Timeout, api.Retries, api.Delay)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
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

func get(rq *http.Request, timeout time.Duration, retries int, retryDelay time.Duration) (*http.Response, error) {
	client := http.Client{
		Timeout: timeout,
	}

	attempts := 0

	var response *http.Response
	var err error

	for {
		attempts += 1
		response, err = client.Do(rq)

		if err == nil {
			if response.StatusCode == http.StatusOK {
				break
			} else {
				err = fmt.Errorf("error getting contact list (%v)", response.Status)
			}
		}

		if attempts >= retries {
			return nil, err
		}

		warn(err)
		time.Sleep(retryDelay)
	}

	return response, nil
}

func info(f string, args ...any) {
	format := fmt.Sprintf("%-5s %v", "INFO", f)

	log.Printf(format, args...)
}

func warn(err error) {
	log.Printf("%-5s %v", "WARN", err)
}
