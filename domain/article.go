package domain

import (
	"fmt"
	"math"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/hyperjumptech/beda"
	levenshtein "github.com/ka-weihe/fast-levenshtein"
)

var N_RELATED_TEXT_SENTENCE int = 2
var MIN_SENTENCE_RELATED_SIZE int = 15

type Article struct {
	Title       string `json:"title"`
	Url         string
	Year        int    `json:"annee"`
	Author      string `json:"author"`
	ID          string `json:"idproprio"`
	Journal     string `json:"titrerev"`
	Text        string `json:"text"`
	RelatedText RelatedText
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

func (a *Article) BuildRelatedText(query string) {
	sentenceSlice := strings.Split(a.Text, ".")

	minScoreIndex := a.findMostRelatedSentenceTrigram(query, sentenceSlice)

	a.RelatedText = a.createdRelatedTextObject(minScoreIndex, sentenceSlice)
}

func (a Article) findMostRelatedSentenceLevenshtein(query string, sentenceSlice []string) int {
	minScore := math.MaxInt
	minScoreIndex := -1
	for i, s := range sentenceSlice {
		if nWord := strings.Split(s, " "); len(nWord) >= MIN_SENTENCE_RELATED_SIZE {
			distance := levenshtein.Distance(query, s)
			if distance < minScore {
				minScore = distance
				minScoreIndex = i
			}
		}

	}
	return minScoreIndex
}

func (a Article) findMostRelatedSentenceTrigram(query string, sentenceSlice []string) int {
	var maxScore float32 = -math.MaxFloat32
	maxScoreIndex := -1
	for i, s := range sentenceSlice {
		if nWord := strings.Split(s, " "); len(nWord) >= MIN_SENTENCE_RELATED_SIZE {
			diff := beda.TrigramCompare(query, s)
			if diff > maxScore {
				maxScore = diff
				maxScoreIndex = i
			}
		}

	}
	return maxScoreIndex
}

func (a Article) createdRelatedTextObject(bestScoreIndex int, sentenceSlice []string) RelatedText {
	lowerBoundRelatedText := bestScoreIndex - N_RELATED_TEXT_SENTENCE
	upperBoundRelatedText := bestScoreIndex + N_RELATED_TEXT_SENTENCE

	if upperBoundRelatedText > len(sentenceSlice)-1 {
		upperBoundRelatedText = len(sentenceSlice) - 1
	}
	if bestScoreIndex != 0 {

		if lowerBoundRelatedText < 0 {
			lowerBoundRelatedText = 0
		}
		relatedTextBefore := ""
		if res := strings.Join(sentenceSlice[lowerBoundRelatedText:bestScoreIndex-1], "."); res != "" {
			relatedTextBefore += res + "."
		}

		relatedTextAfter := ""

		if bestScoreIndex != upperBoundRelatedText {
			relatedTextAfter = strings.Join(sentenceSlice[bestScoreIndex+1:upperBoundRelatedText], ".")
		}

		return RelatedText{
			Prev:  relatedTextBefore,
			Best:  sentenceSlice[bestScoreIndex] + ".",
			After: relatedTextAfter + ".",
		}
	}

	relatedTextAfter := ""

	if bestScoreIndex != upperBoundRelatedText {
		relatedTextAfter = strings.Join(sentenceSlice[1:upperBoundRelatedText], ".")
	}

	return RelatedText{
		Prev:  "",
		Best:  sentenceSlice[bestScoreIndex] + ".",
		After: relatedTextAfter,
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
