package wildapricot

import ()

type Group struct {
	ID   uint32
	Name string
}

type group struct {
	ID          uint32 `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	URL         string `json:"Url"`
	Contacts    int    `json:"ContactsCount"`
}
