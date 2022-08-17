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
		"static/result/header.html",
		"static/result/element_with_description.html",
		"static/result/element_with_persona.html",
		"static/component/input_form.html",
	))

	articles := resp.Data

	pageType := Page{}
	if r.URL.Path == ENTENDU_EN_VOYAGE {
		pageType = Page{
			ResultSectionClass:  "",
			IsEntenduEnVoyage:   true,
			IsRencontreEnVoyage: false,
		}
	} else {
		pageType = Page{
			ResultSectionClass:  "result-grid",
			IsEntenduEnVoyage:   false,
			IsRencontreEnVoyage: true,
		}
	}
	result_info := ResultInfo{
		Results:  articles,
		Query:    resp.Query,
		PageType: pageType,
		NResult:  resp.N,
	}
	err = tmpl.Execute(w, result_info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type ResultInfo struct {
	Results  []domain.Article
	Query    string
	PageType Page
	NResult  int
}

type Page struct {
	ResultSectionClass  string
	IsEntenduEnVoyage   bool
	IsRencontreEnVoyage bool
}
