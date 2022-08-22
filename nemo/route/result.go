package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/erudit-recommandation/search-engine-webapp/config"
	"github.com/erudit-recommandation/search-engine-webapp/domain"
	"github.com/erudit-recommandation/search-engine-webapp/middleware"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const MAX_RESULTS = 10

func Result(w http.ResponseWriter, r *http.Request) {

	var resp middleware.ResultResponse

	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"static/result/results_page.html",
		"static/result/header.html",
		"static/result/element_with_description.html",
		"static/result/element_with_persona.html",
		"static/component/input_form.html",
	))

	articles := resp.Data
	page := managePageType(r, &resp)

	p := message.NewPrinter(language.CanadianFrench)
	articleHashedQuery := make([]ArticleHashedQuery, len(articles))
	for i, a := range resp.Data {
		var personaImageLink string
		if a.PersonaSvg == "" {
			personaImageLink = "/static/images/persona_placeholder.svg"
		} else {
			personaImageLink = fmt.Sprintf("/static/images/persona/%v/%v.svg", resp.HashedQuery, a.ID)
		}
		authorStyle := "none"
		if a.Author == "" {
			a.Author = "Auteur indisponible"
			authorStyle = "italic"
		}

		articleHashedQuery[i] = ArticleHashedQuery{
			Article:          a,
			PersonaImageLink: personaImageLink,
			AuthorStyle:      authorStyle,
			Corpus:           resp.Corpus,
		}

	}
	hostArticleAuthorStyle := "none"
	if resp.HostArticle.Author == "" {
		hostArticleAuthorStyle = "italic"
	}
	result_info := ResultInfo{
		Results:                articleHashedQuery,
		Query:                  resp.Query,
		Page:                   page,
		NResult:                p.Sprintf("%d\n", resp.N),
		HashedQuery:            fmt.Sprintf("%v", resp.HashedQuery),
		Corpus:                 resp.Corpus,
		HostArticle:            resp.HostArticle,
		HostArticleAuthorStyle: hostArticleAuthorStyle,
		AllCorpus:              config.GetConfig().GetCorpusNames()[1:],
	}
	err = tmpl.Execute(w, result_info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func managePageType(r *http.Request, resp *middleware.ResultResponse) Page {
	pageType := Page{}
	if r.URL.Path == RENCONTRE_EN_VOYAGE {
		pageType = Page{
			ResultSectionClass:  "result-grid",
			IsRencontreEnVoyage: true,
		}
	} else if strings.Contains(r.URL.Path, ACCOSTE_EN_VOYAGE_ROOT) {
		pageType = Page{
			ResultSectionClass: "result-grid",
			IsAccosteEnVoyage:  true,
		}
	} else {

		pageType = Page{
			ResultSectionClass: "",
			IsEntenduEnVoyage:  true,
		}

		if resp.Page != 0 {
			pageType.HasPreviousPage = true
			pageType.PreviousPage = fmt.Sprintf("%v/%v?page=%v", ENTENDU_EN_VOYAGE, resp.HashedQuery, resp.Page-1)
		}

		if resp.Page != resp.LastPage {
			pageType.HasNextPage = true
			pageType.NextPage = fmt.Sprintf("%v/%v?page=%v", ENTENDU_EN_VOYAGE, resp.HashedQuery, resp.Page+1)
		}
	}

	return pageType
}

type ResultInfo struct {
	Results                []ArticleHashedQuery
	Query                  string
	HashedQuery            string
	Page                   Page
	NResult                string
	Corpus                 string
	HostArticle            domain.Article
	HostArticleAuthorStyle string
	AllCorpus              []string
}

type ArticleHashedQuery struct {
	domain.Article
	PersonaImageLink string
	AuthorStyle      string
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
