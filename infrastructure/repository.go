package domain

import "github.com/erudit-recommandation/search-engine-webapp/domain"

type ArticlesRepository interface {
	GetByTitle(string, uint) (domain.Article, error)
	GetByAuthor(string, uint) ([]domain.Article, error)
}
