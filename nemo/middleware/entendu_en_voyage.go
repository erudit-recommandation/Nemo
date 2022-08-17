package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"reflect"
	"time"

	"github.com/erudit-recommandation/search-engine-webapp/domain"
	"github.com/erudit-recommandation/search-engine-webapp/infrastructure"
)

func EntenduEnVoyage(next httpHandlerFunc) httpHandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		CACHE.ClearExpired()
		repo, err := infrastructure.ProvideArangoArticlesRepository()
		if err != nil {
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
		log.Printf("-- Entendu en voyage Query: %v --\n", query)
		var resp []domain.Article
		var nFound int
		if r, ok := CACHE[query]; ok {
			nFound = len(r.Elements)
			resp, err = GetEntenduEnvoyageArticleFromCache(query, 0)
			if err != nil {
				Error(w, req, http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			ids, err := repo.GetSearchSentencesID(query, LIMIT_ENTENDU_EN_VOYAGE)
			nFound = len(ids)
			if err != nil {
				Error(w, req, http.StatusInternalServerError, err.Error())
				return
			}
			CACHE[query] = cacheElement{
				CreatedDate: time.Now(),
				Elements:    ids,
			}

			resp, err = GetEntenduEnvoyageArticleFromCache(query, 0)
			if err != nil {
				Error(w, req, http.StatusInternalServerError, err.Error())
				return
			}
		}

		j, err := json.Marshal(ResultResponse{Data: resp, Query: query, N: nFound})

		if err != nil {
			Error(w, req, http.StatusInternalServerError, err.Error())
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(j))
		next(w, req)
	}
}

func GetEntenduEnvoyageArticleFromCache(query string, page uint) ([]domain.Article, error) {
	repo, err := infrastructure.ProvideArangoArticlesRepository()
	if err != nil {
		return nil, err
	}

	resp := make([]domain.Article, 0, MAX_PAGE_ENTENDU_EN_VOYAGE)

	pageIds, err := CACHE[query].GetPage(page)
	if err != nil {
		return nil, err
	}

	for _, id := range pageIds {
		article, err := repo.GetArticleFromSentenceID(id)
		if err != nil {
			return nil, err
		}
		resp = append(resp, article)
	}

	for i := 0; i < len(resp); i++ {
		resp[i].BuildRelatedText(query)
	}
	return resp, nil
}

type cache map[string]cacheElement

func (c *cache) ClearExpired() {
	keys := reflect.ValueOf(*c).MapKeys()
	for _, k := range keys {
		if (*c)[k.String()].IsExpired() {
			delete(*c, k.String())
		}
	}
}

type cacheElement struct {
	CreatedDate time.Time
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
		return nil, fmt.Errorf("this page don't exist")
	}

	if page == c.NumberOfPage() {
		return c.Elements[MAX_PAGE_ENTENDU_EN_VOYAGE*page:], nil
	}

	return c.Elements[MAX_PAGE_ENTENDU_EN_VOYAGE*page : MAX_PAGE_ENTENDU_EN_VOYAGE*(page+1)], nil
}
