package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/erudit-recommandation/search-engine-webapp/config"
	"github.com/erudit-recommandation/search-engine-webapp/domain"
	"github.com/erudit-recommandation/search-engine-webapp/infrastructure"
	"github.com/gorilla/mux"
)

func AccosteEnVoyage(next httpHandlerFunc) httpHandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		corpus := vars["corpus"]
		idproprio := vars["idproprio"]
		hashedQuery := hash(idproprio, corpus)

		log.Printf("-- Accosté en voyage Query: %v --\n", idproprio)

		var nFound int
		var articles []domain.Article
		repo, err := infrastructure.ProvideArangoArticlesRepository(corpus)
		if err != nil {
			log.Println(err)
			Error(w, req, http.StatusInternalServerError, err.Error())
			return
		}

		hostArticle, err := repo.GetByIdproprio(idproprio)
		if err != nil {
			log.Println(err)
			Error(w, req, http.StatusInternalServerError, err.Error())
			return
		}

		if _, ok := CACHE_ACCOSTE_EN_VOYAGE[hashedQuery]; ok {
			resp, err := GetArticleFromCache(hashedQuery, LIMIT_ACCOSTER_EN_VOYAGE, &CACHE_ACCOSTE_EN_VOYAGE)
			if err != nil {
				log.Println(err)
				Error(w, req, http.StatusInternalServerError, err.Error())
				return
			}
			articles = resp
			nFound = len(articles)

		} else {
			corpusInfo, err := config.GetConfig().GetDatabaseCorpus(corpus)
			if err != nil {
				log.Println(err)
				Error(w, req, http.StatusInternalServerError, err.Error())
				return
			}
			resp, err := repo.GetNeighbouringArticlesByBMU(hostArticle.Bmu, corpusInfo.BMUInterval, LIMIT_ACCOSTER_EN_VOYAGE)
			if err != nil {
				log.Println(err)
				Error(w, req, http.StatusInternalServerError, err.Error())
				return
			}
			var anyArticles []interface{} = make([]interface{}, len(resp))
			for i, v := range resp {
				anyArticles[i] = v
			}
			articles = resp
			CACHE_ACCOSTE_EN_VOYAGE[hashedQuery] = newCacheElement("", hashedQuery, anyArticles, LIMIT_ACCOSTER_EN_VOYAGE)
			nFound = len(articles)

			if err := createPersonaSVG(articles, hashedQuery); err != nil {
				log.Println(err)
				Error(w, req, http.StatusInternalServerError, "Il n'y aucun resultat contacter le mainteneur")
				return
			}
		}

		j, err := json.Marshal(ResultResponse{Data: articles, HashedQuery: hashedQuery,
			HostArticle: hostArticle, Corpus: corpus, N: nFound})

		if err != nil {
			log.Println(err)
			Error(w, req, http.StatusInternalServerError, err.Error())
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(j))

		next(w, req)

	}
}
