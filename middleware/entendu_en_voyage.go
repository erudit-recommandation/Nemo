package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/erudit-recommandation/search-engine-webapp/domain"
	"github.com/erudit-recommandation/search-engine-webapp/infrastructure"
)

func EntenduEnVoyage(next httpHandlerFunc) httpHandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		repo, err := infrastructure.ProvideArangoArticlesRepository()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Fprintf(w, "")
			return
		}

		if err := req.ParseForm(); err != nil || req.FormValue("text") == "" {
			err_msg := domain.NO_TEXT_SENDED_FOR_RECOMMANDATION
			log.Println(err)
			http.Error(w, err_msg, http.StatusInternalServerError)
			fmt.Fprintf(w, "")
			return
		}

		query := req.FormValue("text")

		resp, err := repo.SearchPhrases(query, LIMIT)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Fprintf(w, "")
			return
		}

		j, err := json.Marshal(resp)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Fprintf(w, "")
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(j))
		next(w, req)
	}
}
