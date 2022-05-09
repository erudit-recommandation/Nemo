package route

import (
	"net/http"
	"text/template"
)

func Result(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles(
		"static/result/result_page.html",
	))

	err := tmpl.Execute(w, map[string]string{})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

type ResultInfo struct {
	Text string
}
