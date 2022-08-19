package middleware

import (
	"net/http"
	"text/template"
)

func Error(w http.ResponseWriter, r *http.Request, status int, msg string) {
	tmpl := template.Must(template.ParseFiles(
		"static/error/error.html",
	))
	page := ErrorPage{
		StatusCode: status,
		Msg:        msg,
		Contact:    false,
	}
	if status == http.StatusInternalServerError {
		page.Contact = true
	}
	err := tmpl.Execute(w, page)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

type ErrorPage struct {
	StatusCode int
	Msg        string
	Contact    bool
}
