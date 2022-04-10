package parser

type BlockedUsers struct {
	Users []User `json:"users"`
}

type Bookmarks struct {
	Posts []Post `json:"posts"`
}

type Clap struct {
	Post   Post `json:"post"`
	Amount int  `json:"amount"`
}

type Claps struct {
	Claps []Clap `json:"claps"`
}

type Interests struct {
	Publications []Publication `json:"publications"`
	Tags         []Tag         `json:"tags"`
	Topics       []Topic       `json:"topics"`
	Writers      []User        `json:"writers"`
}

type IP struct {
	Address   string `json:"address"`
	CreatedAt string `json:"createdAt"`
}

type IPs struct {
	IPs []IP `json:"ips"`
}

type Post struct {
	Id          string `json:"id"`
	Url         string `json:"url"`
	Title       string `json:"title"`
	PublishedAt string `json:"publishedAt"`
}

type Publication struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Tag struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Topic struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type User struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Url      string `json:"url"`
}
