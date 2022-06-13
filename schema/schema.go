package schema

type MarkupType string
type GrafType string

const (
	A         MarkupType = "a"
	EM        MarkupType = "em"
	STRONG    MarkupType = "strong"
	BR        MarkupType = "br"
	HIGHLIGHT MarkupType = "highlight"
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
	Meta  string `json:"meta,omitempty"`
	Users []User `json:"users"`
}

type Bookmarks struct {
	Meta  string `json:"meta,omitempty"`
	Posts []Post `json:"posts"`
}

type Clap struct {
	Post   Post `json:"post"`
	Amount int  `json:"amount"`
}

type Claps struct {
	Meta  string `json:"meta,omitempty"`
	Claps []Clap `json:"claps"`
}

type Graf struct {
	Type    GrafType `json:"type"`
	Name    string   `json:"name"`
	Text    string   `json:"text,omitempty"`
	Image   *Image   `json:"image,omitempty"`
	Markups []Markup `json:"markups"`
}

type Highlight struct {
	CreatedAt string `json:"createdAt"`
	Body      []Graf `json:"grafs"`
}

type Highlights struct {
	Meta       string      `json:"meta,omitempty"`
	Highlights []Highlight `json:"highlights"`
}

type Image struct {
	Name   string `json:"name,omitempty"`
	Source string `json:"source,omitempty"`
	Alt    string `json:"alt,omitempty"`
	Height string `json:"height,omitempty"`
	Width  string `json:"width,omitempty"`
}

type InnerSection struct {
	Classes []string `json:"classes"`
	Body    []Graf   `json:"body"`
}

type Interests struct {
	Meta         string        `json:"meta,omitempty"`
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
	Meta string `json:"meta,omitempty"`
	IPs  []IP   `json:"ips"`
}

type List struct {
	Name    string `json:"name"`
	Summary string `json:"summary,omitempty"`
	Posts   []Post `json:"posts"`
}

type Lists struct {
	Meta  string `json:"meta,omitempty"`
	Lists []List `json:"list"`
}

type Markup struct {
	Type  MarkupType `json:"type"`
	Start int        `json:"start"`
	End   int        `json:"end"`
	Href  string     `json:"href,omitempty"`
}

type Membership struct {
	Id        string  `json:"id"`
	StartedAt string  `json:"startedAt,omitempty"`
	EndedAt   string  `json:"endedAt,omitempty"`
	Amount    float64 `json:"amount"`
	Type      string  `json:"type,omitempty"`
}

type MembershipCharge struct {
	CreatedAt string  `json:"createdAt"`
	Amount    float64 `json:"amount"`
}

type Post struct {
	Id          string    `json:"id"`
	Url         string    `json:"url"`
	Title       string    `json:"title"`
	PublishedAt string    `json:"publishedAt,omitempty"`
	Content     []Section `json:"content,omitempty"`
}

type Profile struct {
	Meta              string                   `json:"meta,omitempty"`
	User              *User                    `json:"user"`
	Email             string                   `json:"email,omitempty"`
	PastEmails        []string                 `json:"pastEmails,omitempty"`
	MembershipCharges []MembershipCharge       `json:"membershipCharges,omitempty"`
	Memberships       []Membership             `json:"memberships,omitempty"`
	SocialAccounts    map[string]SocialAccount `json:"socialAccounts,omitempty"`
	Editor            []Publication            `json:"editor,omitempty"`
	Writer            []Publication            `json:"writer,omitempty"`
}

type Publication struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Publications struct {
	Meta         string        `json:"meta,omitempty"`
	Publications []Publication `json:"publications"`
}

type Section struct {
	Name string         `json:"name"`
	Body []InnerSection `json:"body"`
}

type Session struct {
	CreatedAt        string `json:"createdAt,omitempty"`
	LastSeenAt       string `json:"lastSeenAt,omitempty"`
	LastSeenLocation string `json:"lastSeenLocation,omitempty"`
	UserAgent        string `json:"userAgent,omitempty"`
}

type Sessions struct {
	Meta     string    `json:"meta,omitempty"`
	Sessions []Session `json:"sessions"`
}

type SocialAccount struct {
	Id    string `json:"id,omitempty"`
	Url   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

type Tag struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Topic struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Topics struct {
	Meta   string  `json:"meta,omitempty"`
	Topics []Topic `json:"topics"`
}

type User struct {
	CreatedAt  string `json:"createdAt,omitempty"`
	Id         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Bio        string `json:"bio,omitempty"`
	Username   string `json:"username"`
	Url        string `json:"url"`
	ProfilePic *Image `json:"profilePic,omitempty"`
}

type Users struct {
	Meta  string `json:"meta,omitempty"`
	Users []User `json:"users"`
}
