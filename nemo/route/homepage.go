package route

import (
	"net/http"
	"text/template"

	"github.com/erudit-recommandation/search-engine-webapp/middleware"
)

func Homepage(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles(
		"static/homepage/homepage.html",
		"static/component/input_form.html",
	))
	err := tmpl.Execute(w, map[string]string{"Query": ""})
	if err != nil {
		middleware.Error(w, r, http.StatusInternalServerError, err.Error())
	}
}
