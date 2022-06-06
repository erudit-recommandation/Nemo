package middleware

import "github.com/erudit-recommandation/search-engine-webapp/domain"

type ResultResponse struct {
	Data  []domain.Article `json:"data"`
	Query string           `json:"query"`
}
