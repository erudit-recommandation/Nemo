package infrastructure

import "github.com/erudit-recommandation/search-engine-webapp/domain"

type ArticlesRepository interface {
	GetByIdproprio(id string) (domain.Article, error)
	SearchPhrases(phrase string, limit uint) ([]domain.Article, error)
}
