package schema

type MarkupType string
type GrafType string

const (
	A      MarkupType = "a"
	EM     MarkupType = "em"
	STRONG MarkupType = "strong"
	BR     MarkupType = "br"
)

const (
	H1         GrafType = "h1"
	H2         GrafType = "h2"
	H3         GrafType = "h3"
	H4         GrafType = "h4"
	IMG        GrafType = "img"
	P          GrafType = "p"
	HR         GrafType = "hr"
	BLOCKQUOTE GrafType = "bq"
	EMBED      GrafType = "embed"
	PRE        GrafType = "pre"
)

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

type Graf struct {
	Type    GrafType `json:"type"`
	Name    string   `json:"name"`
	Text    string   `json:"text,omitempty"`
	Image   *Image   `json:"image,omitempty"`
	Markups []Markup `json:"markups"`
}

type Image struct {
	Name   string `json:"name"`
	Source string `json:"source"`
	Alt    string `json:"alt"`
	Height string `json:"height"`
	Width  string `json:"width"`
}

type InnerSection struct {
	Classes []string `json:"classes"`
	Body    []Graf   `json:"body"`
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

type List struct {
	Name    string `json:"name"`
	Summary string `json:"summary,omitempty"`
	Posts   []Post `json:"posts"`
}

type Lists struct {
	Lists []List `json:"list"`
}

type Markup struct {
	Type  MarkupType `json:"type"`
	Start int        `json:"start"`
	End   int        `json:"end"`
	Href  string     `json:"href,omitempty"`
}

type Post struct {
	Id          string    `json:"id"`
	Url         string    `json:"url"`
	Title       string    `json:"title"`
	PublishedAt string    `json:"publishedAt,omitempty"`
	Content     []Section `json:"content,omitempty"`
}

type Publication struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Section struct {
	Name string         `json:"name"`
	Body []InnerSection `json:"body"`
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
