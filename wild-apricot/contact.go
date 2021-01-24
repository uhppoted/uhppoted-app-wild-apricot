package wildapricot

import (
	"time"
)

type Contact struct {
	ID    uint32
	Name  string
	Email string
}

type contact struct {
	ID                uint32    `json:"Id"`
	Email             string    `json:"Email"`
	FirstName         string    `json:"FirstName"`
	LastName          string    `json:"LastName"`
	DisplayName       string    `json:"DisplayName"`
	Status            string    `json:"Status"`
	MembershipEnabled bool      `json:"MembershipEnabled"`
	Updated           time.Time `json:"ProfileLastUpdated"`
	Fields            []field   `json:"FieldValues"`

	Administrator      bool       `json:"IsAccountAdministrator"`
	MembershipLevel    membership `json:"MembershipLevel"`
	Organization       string     `json:"Organization"`
	TermsOfUseAccepted bool       `json:"TermsOfUseAccepted"`
	URL                string     `json:"Url"`
}

type membership struct {
	ID   uint32 `json:"Id"`
	Name string `json:"Name"`
	URL  string `json:"Url"`
}

type field struct {
	Name       string      `json:"FieldName"`
	SystemCode string      `json:"SystemCode"`
	Value      interface{} `json:"Value"`
}
