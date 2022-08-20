package domain

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
)

var N_RELATED_TEXT_SENTENCE int = 2
var MIN_SENTENCE_RELATED_SIZE int = 15

type Article struct {
	Title            string `json:"title"`
	Url              string
	Year             int    `json:"annee"`
	Author           string `json:"author"`
	ID               string `json:"idproprio"`
	Journal          string `json:"titrerev"`
	CurrentSentence  string `json:"current_sentence"`
	PreviousSentence string `json:"previous_sentence"`
	NextSentence     string `json:"next_sentence"`
	RelatedText      RelatedText
	PersonaSvg       string `json:"persona_svg"`
	Bmu              int    `json:"bmu"`
}

type RelatedText struct {
	Prev  string
	Best  string
	After string
}

func (a *Article) BuildUrl() {
	a.Url = fmt.Sprintf("https://id.erudit.org/iderudit/%v", a.ID)
	if a.Title == "" {
		a.Title = a.Url
	}
}

func (a *Article) BuildRelatedText() {

	a.RelatedText = RelatedText{
		Prev:  a.PreviousSentence,
		Best:  a.CurrentSentence,
		After: a.NextSentence,
	}

}

func NewDummyResults(n int) []Article {
	res := make([]Article, n)
	for i := 0; i < n; i++ {
		res[i] = Article{
			Title:  randomdata.SillyName(),
			Author: randomdata.SillyName(),
			Url:    fmt.Sprintf("https://%v.com", randomdata.SillyName()),
			Year:   randomdata.Number(1900, 2022),
		}
	}
	return res
}
