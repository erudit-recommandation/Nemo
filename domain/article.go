package domain

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
)

type Article struct {
	Title       string `json:"title"`
	Url         string
	Year        int    `json:"annee"`
	Author      string `json:"author"`
	ID          string `json:"idproprio"`
	Journal     string `json:"titrerev"`
	Text        string `json:"text"`
	RelatedText string
}

func (a *Article) BuildUrl() {
	a.Url = fmt.Sprintf("https://id.erudit.org/iderudit/%v", a.ID)
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
