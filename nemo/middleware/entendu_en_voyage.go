package middleware

import (
	"bytes"
	"encoding/json"
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
			Error(w, req, http.StatusInternalServerError, err.Error())
			return
		}

		if err := req.ParseForm(); err != nil || req.FormValue("text") == "" {
			err_msg := domain.NO_TEXT_SENDED_FOR_RECOMMANDATION
			log.Println(err)
			Error(w, req, http.StatusInternalServerError, err_msg)
			return
		}

		query := req.FormValue("text")
		log.Printf("-- Entendu en voyage Query: %v --\n", query)

		resp, err := repo.SearchPhrases(query, LIMIT)

		if err != nil {
			Error(w, req, http.StatusInternalServerError, err.Error())
			return
		}
		for i := 0; i < len(resp); i++ {
			resp[i].BuildRelatedText(query)
		}

		j, err := json.Marshal(ResultResponse{Data: resp, Query: query})

		if err != nil {
			Error(w, req, http.StatusInternalServerError, err.Error())
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(j))
		next(w, req)
	}
}
