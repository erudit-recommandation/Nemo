package middleware

import "github.com/erudit-recommandation/search-engine-webapp/domain"

type ResultResponse struct {
	Data        []domain.Article `json:"data"`
	Query       string           `json:"query"`
	HashedQuery uint32           `json:"hased_query"`
	N           int              `json:"n"`
	Page        uint             `json:"page"`
	LastPage    uint             `json:"last_page"`
	HostArticle domain.Article   `json:"host_article"`
	Corpus      string           `json:"corpus"`
}
