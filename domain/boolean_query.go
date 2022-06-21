package domain

import (
	"regexp"
	"strings"
)

type Operation int

const (
	AND Operation = iota
	OR
	NOT
)

type BooleanQuery struct {
	phrases    []string
	operations []Operation
}

func (b BooleanQuery) Phrase() []string {
	return b.phrases
}

func (b BooleanQuery) Operations() []Operation {
	return b.operations
}

func NewBooleanQuery(query string) BooleanQuery {
	re := regexp.MustCompile("(AND|OR|NOT)")
	if re.MatchString(query) { //boolean operation
		operations := associateOperation(re.FindAllString(query, -1))
		phrases := re.Split(query, -1)

		return BooleanQuery{
			phrases:    phrases,
			operations: operations,
		}

	} else if sentences := strings.Split(query, "."); len(sentences) > 1 { // multiple phrases
		operations := make([]Operation, len(sentences)-1)
		for i := 0; i < len(sentences)-1; i += 1 {
			operations[i] = AND
		}
		return BooleanQuery{
			phrases:    sentences,
			operations: operations,
		}
	}
	// normal case with no operation
	return BooleanQuery{
		phrases: []string{query},
	}
}

func associateOperation(operations []string) []Operation {
	resp := make([]Operation, len(operations))
	for i, o := range operations {
		if o == "AND" {
			resp[i] = AND
		} else if o == "OR" {
			resp[i] = OR
		} else {
			resp[i] = NOT
		}
	}

	return resp
}

// abc. def.
// phrase("abc. def.", text_fr)
// phrase (abc) AND phrase(def)
