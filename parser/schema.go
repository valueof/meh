package parser

type BlockedList struct {
	Users []MediumUser `json:"users"`
}

type MediumUser struct {
	Username string `json:"username"`
}
