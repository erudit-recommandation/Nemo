package domain

type ArticlesRepository interface {
	GetByTitle(string, uint) (Article, error)
	GetByAuthor(string, uint) ([]Article, error)
}
