package wildapricot

import ()

type MemberGroup struct {
	ID          uint32 `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	URL         string `json:"Url"`
	Contacts    int    `json:"ContactsCount"`
}
