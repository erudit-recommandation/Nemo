package route

import (
	"net/http"
	"text/template"

	"github.com/erudit-recommandation/search-engine-webapp/middleware"
)

func Vaiseau(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles(
		"static/vaiseau/vaiseau.html",
	))
	err := tmpl.Execute(w, map[string]string{})
	if err != nil {
		middleware.Error(w, r, http.StatusInternalServerError, err.Error())
	}
}
