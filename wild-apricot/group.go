package wildapricot

import ()

type MemberGroup struct {
	ID          uint32 `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	URL         string `json:"Url"`
	Contacts    int    `json:"ContactsCount"`
}

func (mg *MemberGroup) Flatten() (map[string]interface{}, error) {
	flattened := map[string]interface{}{}

	if mg != nil {
		flattened["id"] = mg.ID
		flattened["name"] = mg.Name
		flattened["description"] = mg.Description
		flattened["contacts"] = mg.Contacts
		flattened["url"] = mg.URL
	}

	return flattened, nil
}
