package route

import (
	"net/http"
	"text/template"

	"github.com/erudit-recommandation/search-engine-webapp/config"
	"github.com/erudit-recommandation/search-engine-webapp/middleware"
)

func Homepage(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles(
		"static/homepage/homepage.html",
		"static/component/input_form.html",
	))
	err := tmpl.Execute(w, homepageInfo{
		Query:         "",
		AllCorpus:     config.GetConfig().GetCorpusNames()[1:],
		DefaultCorpus: config.GetConfig().GetCorpusNames()[0],
	})

	if err != nil {
		middleware.Error(w, r, http.StatusInternalServerError, err.Error())
	}
}

type homepageInfo struct {
	Query         string
	AllCorpus     []string
	DefaultCorpus string
	ResultInfo
}
