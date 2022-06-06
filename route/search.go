package route

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/erudit-recommandation/search-engine-webapp/domain"
)

const MAX_RESULTS = 10

func Result(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil || r.FormValue("text") == "" {
		err_msg := domain.NO_TEXT_SENDED_FOR_RECOMMANDATION
		log.Println(err)
		http.Error(w, err_msg, http.StatusInternalServerError)
		fmt.Fprintf(w, "")
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"static/result/results_page.html",
		"static/result/result.html",
	))

	result_info := ResultInfo{
		Results: domain.NewDummyResults(MAX_RESULTS),
	}
	err := tmpl.Execute(w, result_info)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

type ResultInfo struct {
	Results []domain.Article
}
