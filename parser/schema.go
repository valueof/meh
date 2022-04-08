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

type Post struct {
	Id          string `json:"id"`
	Url         string `json:"url"`
	Title       string `json:"title"`
	PublishedAt string `json:"publishedAt"`
}

type User struct {
	Username string `json:"username"`
}
