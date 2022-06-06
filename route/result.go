package route

import (
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/erudit-recommandation/search-engine-webapp/domain"
)

const MAX_RESULTS = 10

func Result(w http.ResponseWriter, r *http.Request) {

	var articles []domain.Article

	err := json.NewDecoder(r.Body).Decode(&articles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"static/result/results_page.html",
		"static/result/result.html",
	))

	result_info := ResultInfo{
		Results: articles,
	}
	err = tmpl.Execute(w, result_info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type ResultInfo struct {
	Results []domain.Article
}
