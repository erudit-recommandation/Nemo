package domain

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
)

type Article struct {
	Title    string `json:"title"`
	Abstract string
	Url      string
	Date     string `json:"annee"`
	DOI      string
}

func NewDummyResults(n int) []Article {
	res := make([]Article, n)
	for i := 0; i < n; i++ {
		res[i] = Article{
			Title:    randomdata.SillyName(),
			Abstract: randomdata.Paragraph(),
			Url:      fmt.Sprintf("https://%v.com", randomdata.SillyName()),
			Date:     fmt.Sprintf("%v/%v/%v", randomdata.Number(0, 31), randomdata.Number(0, 12), randomdata.Number(1900, 2022)),
			DOI:      fmt.Sprintf("https://doi/%v.com", randomdata.SillyName()),
		}
	}
	return res
}
