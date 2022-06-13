package parser

var images map[string]bool

func init() {
	images = map[string]bool{}
}

func GetImages() map[string]bool {
	return images
}
