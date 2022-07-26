package server

import (
	"net/http"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	data := pageMeta{}
	data.Title = "[meh] Page Not Found"
	data.SkipFooter = true

	render(w, r, "404.html", data)
}

type internalServerErrorData struct {
	RequestID string
	pageMeta
}

func internalServerError(w http.ResponseWriter, r *http.Request) {
	data := internalServerErrorData{}
	data.Title = "[meh] Internal Server Error"
	data.SkipFooter = true
	data.RequestID = getRequestIDFromContext(r.Context())

	render(w, r, "500.html", data)
}
