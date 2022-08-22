package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/erudit-recommandation/search-engine-webapp/config"
	"github.com/erudit-recommandation/search-engine-webapp/domain"
	"github.com/erudit-recommandation/search-engine-webapp/infrastructure"
)

func RencontreEnVoyage(next httpHandlerFunc) httpHandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		corpus := "erudit"
		if err := req.ParseForm(); err != nil || req.FormValue("text") == "" {
			err_msg := domain.NO_TEXT_SENDED_FOR_RECOMMANDATION
			log.Println(err)

			Error(w, req, http.StatusBadRequest, err_msg)
			return
		}

		query := req.FormValue("text")
		log.Printf("-- Rencontré en voyage Query: %v --\n", query)
		hashedQuery := hash(query, corpus)

		var articles []domain.Article

		if _, ok := CACHE_RENCONTRE_EN_VOYAGE[hashedQuery]; ok {
			resp, err := GetArticleFromCache(hashedQuery, LIMIT_RENCONTRE_EN_VOYAGE, &CACHE_RENCONTRE_EN_VOYAGE)
			articles = resp
			if err != nil {
				log.Println(err)
				Error(w, req, http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			recommandation, err := sendRequestToGemsimService(query, LIMIT_RENCONTRE_EN_VOYAGE, corpus)
			if err != nil {
				log.Println(err)
				Error(w, req, http.StatusInternalServerError, err.Error())
				return
			}
			repo, err := infrastructure.ProvideArangoArticlesRepository(corpus)
			if err != nil {
				log.Println(err)
				Error(w, req, http.StatusInternalServerError, err.Error())
				fmt.Fprintf(w, "")
				return
			}
			articlesParsed := make([]articleScore, 0, LIMIT_RENCONTRE_EN_VOYAGE)

			for k, v := range recommandation {
				id, err := strconv.Atoi(k)
				if err != nil {
					log.Println(err)
					Error(w, req, http.StatusInternalServerError, err.Error())
					return
				}
				article, err := repo.GetByIdPandas(id)
				if err == nil {
					articlesParsed = append(articlesParsed, articleScore{Article: article, Score: v})
					sort.Slice(articlesParsed, func(i, j int) bool { return articlesParsed[i].Score > articlesParsed[j].Score })

				}

			}

			articles = make([]domain.Article, 0, len(articlesParsed))

			for _, el := range articlesParsed {
				articles = append(articles, el.Article)
			}

			if len(articles) == 0 {
				log.Println(err)
				Error(w, req, http.StatusInternalServerError, "Il n'y aucun resultat contacter le mainteneur")
				return
			}

			var anyArticles []interface{} = make([]interface{}, len(articles))
			for i, v := range articles {
				anyArticles[i] = v
			}
			CACHE_RENCONTRE_EN_VOYAGE[hashedQuery] = newCacheElement(query, hashedQuery, anyArticles, LIMIT_RENCONTRE_EN_VOYAGE)

			if err := createPersonaSVG(articles, hashedQuery); err != nil {
				log.Println(err)
				Error(w, req, http.StatusInternalServerError, "Il n'y aucun resultat contacter le mainteneur")
				return
			}

		}

		j, err := json.Marshal(ResultResponse{Data: articles, Query: query, HashedQuery: hashedQuery, Corpus: corpus})

		if err != nil {
			log.Println(err)
			Error(w, req, http.StatusInternalServerError, err.Error())
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(j))

		next(w, req)
	}
}

func sendRequestToGemsimService(text string, n uint, corpus string) (map[string]float64, error) {
	gemsimAddr := fmt.Sprintf("%v/gensim", config.GetConfig().TextAnalysisServiceAddr)
	body := gemsimServiceRequest{
		Text:   text,
		N:      n,
		Corpus: corpus,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(gemsimAddr, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	var responseMap map[string]float64

	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return nil, err
	}
	return responseMap, nil
}

type articleScore struct {
	Article domain.Article
	Score   float64
}
