package parser

type BlockedUsers struct {
	Users []User `json:"users"`
}

type Bookmarks struct {
	Posts []Post `json:"posts"`
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
