package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/erudit-recommandation/search-engine-webapp/domain"
	"github.com/erudit-recommandation/search-engine-webapp/infrastructure"
	"github.com/gorilla/mux"
)

func EntenduEnVoyage(next httpHandlerFunc) httpHandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		CACHE.ClearExpired()
		repo, err := infrastructure.ProvideArangoArticlesRepository()
		if err != nil {
			log.Println(err)
			Error(w, req, http.StatusInternalServerError, err.Error())
			return
		}

		if err := req.ParseForm(); err != nil || req.FormValue("text") == "" {
			err_msg := domain.NO_TEXT_SENDED_FOR_RECOMMANDATION
			log.Println(err)
			Error(w, req, http.StatusBadRequest, err_msg)
			return
		}

		query := req.FormValue("text")

		page := 0

		log.Printf("-- Entendu en voyage Query: %v --\n", query)

		var resp []domain.Article
		var nFound int
		var lastPage uint
		hasedQuery := hash(query)

		if r, ok := CACHE[hasedQuery]; ok {
			lastPage = r.NumberOfPage()
			nFound = len(r.Elements)
			articles, errorCode, err := GetEntenduEnvoyageArticleFromCache(repo, hasedQuery, uint(page))
			resp = articles
			if err != nil {
				Error(w, req, errorCode, err.Error())
				return
			}
		} else {
			ids, err := repo.GetSearchSentencesID(query, LIMIT_ENTENDU_EN_VOYAGE)

			nFound = len(ids)
			if err != nil {
				log.Println(err)
				Error(w, req, http.StatusInternalServerError, err.Error())
				return
			}
			CACHE[hasedQuery] = cacheElement{
				CreatedDate: time.Now(),
				Elements:    ids,
				Query:       query,
			}

			articles, errorCode, err := GetEntenduEnvoyageArticleFromCache(repo, hasedQuery, 0)
			lastPage = CACHE[hasedQuery].NumberOfPage()
			resp = articles
			if err != nil {
				log.Println(err)
				Error(w, req, errorCode, err.Error())
				return
			}
		}

		j, err := json.Marshal(ResultResponse{Data: resp, Query: query, N: nFound, Page: uint(page), LastPage: lastPage, HashedQuery: hasedQuery})

		if err != nil {
			Error(w, req, http.StatusInternalServerError, err.Error())
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(j))
		next(w, req)
	}
}

func EntenduEnVoyageCached(next httpHandlerFunc) httpHandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		CACHE.ClearExpired()

		repo, err := infrastructure.ProvideArangoArticlesRepository()
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

		r, ok := CACHE[uint32(hasedQuery)]
		if !ok {
			err_msg := "Cette requête n'a jamais été faite"
			log.Println(err_msg)
			Error(w, req, http.StatusNotFound, err_msg)
			return
		}
		var page int
		pageString := req.URL.Query().Get("page")
		log.Printf("-- Entendu en voyage Page: %v --\n", pageString)

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

		lastPage := r.NumberOfPage()
		nFound := len(r.Elements)

		articles, errorCode, err := GetEntenduEnvoyageArticleFromCache(repo, uint32(hasedQuery), uint(page))
		resp := articles
		if err != nil {
			Error(w, req, errorCode, err.Error())
			return
		}

		j, err := json.Marshal(ResultResponse{Data: resp, Query: r.Query, N: nFound, Page: uint(page), LastPage: lastPage, HashedQuery: uint32(hasedQuery)})

		if err != nil {
			Error(w, req, http.StatusInternalServerError, err.Error())
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(j))
		next(w, req)

	}
}

func GetEntenduEnvoyageArticleFromCache(repo infrastructure.ArticlesRepository, hasedQuery uint32, page uint) ([]domain.Article, int, error) {

	resp := make([]domain.Article, 0, MAX_PAGE_ENTENDU_EN_VOYAGE)

	pageIds, err := CACHE[hasedQuery].GetPage(page)
	if err != nil {
		return nil, http.StatusNotFound, err
	}

	for _, id := range pageIds {
		article, err := repo.GetArticleFromSentenceID(id)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		resp = append(resp, article)
	}

	for i := 0; i < len(resp); i++ {
		resp[i].BuildRelatedText()
	}
	return resp, http.StatusOK, nil
}

type cache map[uint32]cacheElement

func (c *cache) ClearExpired() {
	keys := reflect.ValueOf(*c).MapKeys()
	for _, k := range keys {
		if (*c)[uint32(k.Uint())].IsExpired() {
			delete(*c, uint32(k.Uint()))
		}
	}
}

type cacheElement struct {
	CreatedDate time.Time
	Query       string
	Elements    []infrastructure.ArticlesID
}

func (c cacheElement) IsExpired() bool {
	return time.Until(c.CreatedDate) >= CACHE_DURATION
}

func (c cacheElement) NumberOfPage() uint {
	return uint(math.Ceil(float64(len(c.Elements)) / float64(MAX_PAGE_ENTENDU_EN_VOYAGE)))
}

func (c cacheElement) GetPage(page uint) ([]infrastructure.ArticlesID, error) {
	if page > c.NumberOfPage() {
		return nil, fmt.Errorf("la page n'existe pas")
	}

	if page == c.NumberOfPage() {
		return c.Elements[MAX_PAGE_ENTENDU_EN_VOYAGE*page:], nil
	}

	return c.Elements[MAX_PAGE_ENTENDU_EN_VOYAGE*page : MAX_PAGE_ENTENDU_EN_VOYAGE*(page+1)], nil
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
