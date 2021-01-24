package wildapricot

import ()

type Contact struct {
	ID    uint32
	Name  string
	Email string
}

type contact struct {
	ID                 uint32     `json:"Id"`
	Email              string     `json:"Email"`
	DisplayName        string     `json:"DisplayName"`
	FirstName          string     `json:"FirstName"`
	LastName           string     `json:"LastName"`
	MembershipEnabled  bool       `json:"MembershipEnabled"`
	Status             string     `json:"Status"`
	MembershipLevel    membership `json:"MembershipLevel"`
	Organization       string     `json:"Organization"`
	TermsOfUseAccepted bool       `json:"TermsOfUseAccepted"`
	Fields             []field    `json:"FieldValues"`
}

type membership struct {
	ID   uint32 `json:"Id"`
	Name string `json:"Name"`
}

type field struct {
	Name       string      `json:"FieldName"`
	SystemCode string      `json:"SystemCode"`
	Value      interface{} `json:"Value"`
}
