package route

import (
	"net/http"
	"text/template"
)

func Homepage(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles(
		"static/homepage/homepage.html",
		"static/component/input_form.html",
	))
	err := tmpl.Execute(w, map[string]string{})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
