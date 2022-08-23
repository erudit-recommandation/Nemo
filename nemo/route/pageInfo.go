package route

import (
	"github.com/erudit-recommandation/search-engine-webapp/config"
	"github.com/erudit-recommandation/search-engine-webapp/domain"
)

type ResultInfo struct {
	Results        []ArticleHashedQuery
	Query          string
	HashedQuery    string
	Page           Page
	NResult        string
	Corpus         config.DatabaseCorpus
	HostArticle    domain.Article
	ResofTheCorpus []config.DatabaseCorpus
}

type ArticleHashedQuery struct {
	domain.Article
	PersonaImageLink string
	Corpus           string
}

type Page struct {
	ResultSectionClass  string
	IsEntenduEnVoyage   bool
	IsRencontreEnVoyage bool
	IsAccosteEnVoyage   bool

	CurrentPage  string
	NextPage     string
	PreviousPage string

	HasNextPage     bool
	HasPreviousPage bool
}
