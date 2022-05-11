package domain

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
)

type Result struct {
	Title    string
	Abstract string
	Url      string
	Date     string
	DOI      string
}

func NewDummyResults(n int) []Result {
	res := make([]Result, n)
	for i := 0; i < n; i++ {
		res[i] = Result{
			Title:    randomdata.SillyName(),
			Abstract: randomdata.Paragraph(),
			Url:      fmt.Sprintf("https://%v.com", randomdata.SillyName()),
			Date:     fmt.Sprintf("%v/%v/%v", randomdata.Number(0, 31), randomdata.Number(0, 12), randomdata.Number(1900, 2022)),
			DOI:      fmt.Sprintf("https://doi/%v.com", randomdata.SillyName()),
		}
	}
	return res
}
