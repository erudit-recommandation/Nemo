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
		n := 20
		query := req.FormValue("text")
		log.Printf("-- Rencontré en voyage Query: %v --\n", query)

		recommandation, err := sendRequestToGemsimService(query, n)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Fprintf(w, "")
			return
		}
		repo, err := infrastructure.ProvideArangoArticlesRepository()
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Fprintf(w, "")
			return
		}
		articlesParsed := make([]articleScore, 0, n)

		for k, v := range recommandation {
			id, err := strconv.Atoi(k)
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				fmt.Fprintf(w, "")
				return
			}
			article, err := repo.GetByIdPandas(id)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				fmt.Fprintf(w, "")
				return
			}
			articlesParsed = append(articlesParsed, articleScore{Article: article, Score: v})
			sort.Slice(articlesParsed, func(i, j int) bool { return articlesParsed[i].Score > articlesParsed[j].Score })
		}

		articles := make([]domain.Article, len(articlesParsed))

		for i, el := range articlesParsed {
			articles[i] = el.Article
		}
		j, err := json.Marshal(ResultResponse{Data: articles, Query: query})

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			fmt.Fprintf(w, "")
			return
		}
		req.Body = ioutil.NopCloser(bytes.NewReader(j))

		next(w, req)
	}
}

func sendRequestToGemsimService(text string, n int) (map[string]float64, error) {
	gemsimAddr := fmt.Sprintf("%v/gemsim", config.GetConfig().GemsimServiceAddr)
	body := gemsimServiceRequest{
		Text: text,
		N:    n,
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
