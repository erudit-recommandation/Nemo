package route

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/erudit-recommandation/search-engine-webapp/domain"
)

func Result(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil || r.FormValue("text") == "" {
		err_msg := domain.NO_TEXT_SENDED_FOR_RECOMMANDATION
		log.Println(err)
		http.Error(w, err_msg, http.StatusInternalServerError)
		fmt.Fprintf(w, "")
		return
	}

	text := r.FormValue("text")
	log.Println(text)

	tmpl := template.Must(template.ParseFiles(
		"static/result/result_page.html",
	))

	result_info := ResultInfo{
		Text: text,
	}
	err := tmpl.Execute(w, result_info)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

type ResultInfo struct {
	Text string
}
