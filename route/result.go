package route

import (
	"encoding/json"
	"net/http"
	"text/template"

	"github.com/erudit-recommandation/search-engine-webapp/domain"
	"github.com/erudit-recommandation/search-engine-webapp/middleware"
)

const MAX_RESULTS = 10

func Result(w http.ResponseWriter, r *http.Request) {

	var resp middleware.ResultResponse

	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"static/result/results_page.html",
		"static/result/result.html",
	))

	articles := resp.Data

	for i := 0; i < len(articles); i++ {
		articles[i].BuildRelatedText(resp.Query)
	}

	result_info := ResultInfo{
		Results: articles,
		Query:   resp.Query,
	}
	err = tmpl.Execute(w, result_info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type ResultInfo struct {
	Results []domain.Article
	Query   string
}
