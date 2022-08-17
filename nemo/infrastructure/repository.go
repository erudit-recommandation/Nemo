package infrastructure

import "github.com/erudit-recommandation/search-engine-webapp/domain"

type ArticlesRepository interface {
	GetByIdproprio(id string) (domain.Article, error)
	GetByIdPandas(id int) (domain.Article, error)
	SearchSentences(phrase string, limit uint) ([]domain.Article, error)

	GetSearchSentencesID(phrase string, limit uint) ([]ArticlesID, error)
	GetArticleFromSentenceID(articleID ArticlesID) (domain.Article, error)
}

type ArticlesID struct {
	Id        string `json:"id"`
	NSentence int    `json:"n_sentence"`
}
