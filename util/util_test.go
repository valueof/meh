package util_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/valueof/meh/schema"
	"github.com/valueof/meh/util"
)

func TestParseMediumId(t *testing.T) {
	tests := map[string]string{
		"https://anton.medium.com/birding-report-july-4th-7e904c599273":              "7e904c599273",
		"https://medium.engineering/simple-style-sheets-c3b588867899":                "c3b588867899",
		"https://medium.com/programming-is-a-nightmare/heaven-and-hell-cb1ec71a9d4a": "cb1ec71a9d4a",
		"https://medium.com/@anton":                                                  "",
	}

	for url, want := range tests {
		have := util.ParseMediumId(url)
		if want != have {
			t.Errorf("want: %s; have: %s", want, have)
		}
	}
}

func TestParseMediumUsername(t *testing.T) {
	tests := map[string]string{
		"https://anton.medium.com/":      "anton",
		"https://medium.com/@anton":      "anton",
		"https://medium.com":             "",
		"https://anton.medium.com/about": "anton",
	}

	for url, want := range tests {
		have := util.ParseMediumUsername(url)
		if want != have {
			t.Errorf("url: %s; want: %s; have: %s", url, want, have)
		}
	}
}

func TestText(t *testing.T) {
	tests := map[string]string{
		//`<p>The <em>owls</em> are not what <strong><em>they seem</em></strong></p>`: "The owls are not what they seem",
		`<p>
			The <em>owls</em>
			are not what
			<strong>
				<em>they seem</em>
			</strong>
		</p>
		`: "The owls are not what they seem",
	}

	for src, want := range tests {
		node, err := util.NewNodeFromHTML(strings.NewReader(src))
		if err != nil {
			t.Errorf("%v", err)
			continue
		}

		have := node.Text()
		if have != want {
			t.Errorf("want: %s; have: %s", want, have)
		}
	}

}

func TestMarkup(t *testing.T) {
	tests := []string{
		`<p>The <em>owls</em> are not what <strong><em>they seem</em></strong></p>`,
		`<p>
			The <em>owls</em>
			are not what
			<strong>
				<em>they seem</em>
			</strong>
		</p>
		`,
	}

	want := []schema.Markup{
		{Type: schema.EM, Start: 4, End: 8},
		{Type: schema.STRONG, Start: 22, End: 31},
		{Type: schema.EM, Start: 22, End: 31},
	}

	for n, tt := range tests {
		node, err := util.NewNodeFromHTML(strings.NewReader(tt))
		if err != nil {
			t.Errorf("test %d failed: %v", n, err)
			continue
		}

		have := node.Markup()
		if reflect.DeepEqual(have, want) == false {
			t.Errorf("test %d failed", n)
			t.Errorf("want: %v; have: %v", want, have)
		}
	}
}
