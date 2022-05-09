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

func NewDummyResults() Result {
	return Result{
		Title:    randomdata.SillyName(),
		Abstract: randomdata.Paragraph(),
		Url:      fmt.Sprintf("https://%v.com", randomdata.SillyName()),
		Date:     fmt.Sprintf("%v/%v/%v", randomdata.Day(), randomdata.Month(), randomdata.Number(1900, 2022)),
		DOI:      fmt.Sprintf("https://%v.com", randomdata.SillyName()),
	}
}
