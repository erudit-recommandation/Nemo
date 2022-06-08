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

func (a Article) createdRelatedTextObject(minScoreIndex int, sentenceSlice []string) RelatedText {

	lowerBoundRelatedText := minScoreIndex - N_RELATED_TEXT_SENTENCE
	upperBoundRelatedText := minScoreIndex + N_RELATED_TEXT_SENTENCE

	if lowerBoundRelatedText < 0 {
		lowerBoundRelatedText = 0
	}

	if upperBoundRelatedText > len(sentenceSlice)-1 {
		upperBoundRelatedText = len(sentenceSlice) - 1
	}
	relatedTextAfter := ""

	if minScoreIndex != upperBoundRelatedText {
		relatedTextAfter = strings.Join(sentenceSlice[minScoreIndex+1:upperBoundRelatedText], ".")
	}

	return RelatedText{
		Prev:  strings.Join(sentenceSlice[lowerBoundRelatedText:minScoreIndex-1], ".") + ".",
		Best:  sentenceSlice[minScoreIndex] + ".",
		After: relatedTextAfter + ".",
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
