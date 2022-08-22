package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/erudit-recommandation/search-engine-webapp/domain"
	"github.com/erudit-recommandation/search-engine-webapp/infrastructure"
	"github.com/gorilla/mux"
)

func EntenduEnVoyage(next httpHandlerFunc) httpHandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		CACHE_ENTENDU_EN_VOYAGE.ClearExpired()

		if err := req.ParseForm(); err != nil || (req.FormValue("text") == "" && req.FormValue("corpus") == "") {
			err_msg := domain.NO_TEXT_SENDED_FOR_RECOMMANDATION
			log.Println(err)
			Error(w, req, http.StatusBadRequest, err_msg)
			return
		}

		query := req.FormValue("text")
		corpus := req.FormValue("corpus")
		repo, err := infrastructure.ProvideArangoArticlesRepository(corpus)
		if err != nil {
			log.Println(err)
			Error(w, req, http.StatusInternalServerError, err.Error())
			return
		}

		page := 0

		log.Printf("-- Entendu en voyage Query: %v --\n", query)

		var resp []domain.Article
		var nFound int
		var lastPage uint
		hashedQuery := hash(query, corpus)

		if r, ok := CACHE_ENTENDU_EN_VOYAGE[hashedQuery]; ok {
			lastPage = uint(r.NumberOfPage())
			nFound = len(r.Elements)
			articles, errorCode, err := GetEntenduEnvoyageArticleFromCache(repo, hashedQuery, uint(page))
			resp = articles
			if err != nil {
				log.Println(err)
				Error(w, req, errorCode, err.Error())
				return
			}
		} else {
			ids, err := repo.GetSearchSentencesID(query, LIMIT_ENTENDU_EN_VOYAGE)
			if err != nil {
				log.Println(err)
				Error(w, req, http.StatusInternalServerError, err.Error())
				return
			}

			nFound = len(ids)
			lastPage = 0
			if nFound != 0 {
				var anyIds []interface{} = make([]interface{}, len(ids))
				for i, v := range ids {
					anyIds[i] = v
				}
				CACHE_ENTENDU_EN_VOYAGE[hashedQuery] = newCacheElement(query, hashedQuery, anyIds, MAX_BY_PAGE_ENTENDU_EN_VOYAGE)

				articles, errorCode, err := GetEntenduEnvoyageArticleFromCache(repo, hashedQuery, 0)
				lastPage = uint(CACHE_ENTENDU_EN_VOYAGE[hashedQuery].NumberOfPage())
				resp = articles
				if err != nil {
					log.Println(err)
					Error(w, req, errorCode, err.Error())
					return
				}
			}

		}

		j, err := json.Marshal(ResultResponse{Data: resp, Query: query,
			N: nFound, Page: uint(page), LastPage: lastPage, HashedQuery: hashedQuery, Corpus: corpus})

		if err != nil {
			log.Println(err)
			Error(w, req, http.StatusInternalServerError, err.Error())
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(j))
		next(w, req)
	}
}

func EntenduEnVoyageCached(next httpHandlerFunc) httpHandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		corpus := "erudit"
		CACHE_ENTENDU_EN_VOYAGE.ClearExpired()

		repo, err := infrastructure.ProvideArangoArticlesRepository(corpus)
		if err != nil {
			log.Println(err)
			Error(w, req, http.StatusInternalServerError, err.Error())
			return
		}
		vars := mux.Vars(req)
		hasedQuery, err := strconv.ParseUint(vars["hashedQuery"], 10, 32)
		if err != nil {
			log.Println(err)
			Error(w, req, http.StatusInternalServerError, err.Error())
			return
		}

		r, ok := CACHE_ENTENDU_EN_VOYAGE[uint32(hasedQuery)]
		if !ok {
			err_msg := "Cette requête n'a jamais été faite"
			log.Println(err_msg)
			Error(w, req, http.StatusNotFound, err_msg)
			return
		}
		var page int
		pageString := req.URL.Query().Get("page")
		log.Printf("-- Entendu en voyage Page: %v; Query: %v --\n", pageString, r.Query)

		if pageString == "" {
			page = 0
		} else {
			page, err = strconv.Atoi(pageString)
			if err != nil {
				err_msg := "la page doit être un nombre"
				log.Println(err)
				Error(w, req, http.StatusBadRequest, err_msg)
				return
			}

		}

		lastPage := uint(r.NumberOfPage())
		nFound := len(r.Elements)

		articles, errorCode, err := GetEntenduEnvoyageArticleFromCache(repo, uint32(hasedQuery), uint(page))
		resp := articles
		if err != nil {
			log.Println(err)
			Error(w, req, errorCode, err.Error())
			return
		}

		j, err := json.Marshal(ResultResponse{Data: resp, Query: r.Query, N: nFound, Page: uint(page), LastPage: lastPage, HashedQuery: uint32(hasedQuery), Corpus: corpus})

		if err != nil {
			log.Println(err)
			Error(w, req, http.StatusInternalServerError, err.Error())
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(j))
		next(w, req)

	}
}

func GetEntenduEnvoyageArticleFromCache(repo infrastructure.ArticlesRepository, hasedQuery uint32, page uint) ([]domain.Article, int, error) {

	resp := make([]domain.Article, 0, MAX_BY_PAGE_ENTENDU_EN_VOYAGE)
	el := CACHE_ENTENDU_EN_VOYAGE[hasedQuery]
	pageIds, err := el.GetPage(page)

	if err != nil {
		return nil, http.StatusNotFound, err
	}
	for _, id := range pageIds {
		article, err := repo.GetArticleFromSentenceID(id.(infrastructure.ArticlesID))
		if err != nil {
			log.Println(pageIds)
			log.Println(el)
			return nil, http.StatusInternalServerError, err
		}
		resp = append(resp, article)
	}

	for i := 0; i < len(resp); i++ {
		resp[i].BuildRelatedText()
	}

	if err := createPersonaSVG(resp, hasedQuery); err != nil {
		return resp, http.StatusInternalServerError, err
	}
	return resp, http.StatusOK, nil
}
