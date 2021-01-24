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

	//	type permission struct {
	//		AccountID         uint32   `json:"AccountId"`
	//		SecurityProfileId uint32   `json:"SecurityProfileId"`
	//		AvailableScopes   []string `json:AvailableScopes"`
	//	}

	result := struct {
		AccessToken string `json:"access_token"`
		//		TokenType    string       `json:"token_type"`
		//		ExpiresIn    int          `json:"expires_in"`
		//		RefreshToken string       `json:"refresh_token"`
		//		Permissions  []permission `json:"Permissions"`
	}{}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
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

	data := struct {
		Contacts []contact `json:"Contacts"`
	}{}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	contacts := []Contact{}
	for _, c := range data.Contacts {
		contacts = append(contacts, Contact{
			ID:    c.ID,
			Name:  c.DisplayName,
			Email: c.Email,
		})
	}

	return contacts, nil
}
