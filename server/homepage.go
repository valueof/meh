package server

import (
	"net/http"
)

func homepage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFound(w, r)
		return
	}

	render(w, r, "home.html", pageMeta{Title: "Medium Export Helper"})
}
