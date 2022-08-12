package domain

import (
	"fmt"
	"regexp"
	"strings"
)

type Operation int

const (
	AND Operation = iota
	OR
	NOT
)

func (o Operation) String() string {
	switch o {
	case AND:
		return "AND"
	case OR:
		return "OR"
	case NOT:
		return "AND NOT"
	}
	return ""
}

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

func (b BooleanQuery) ToArangoPhraseQueryBody() string {
	resp := ""
	resp += fmt.Sprintf(`CONTAINS(LOWER(s.text), LOWER("%v")) `, b.phrases[0])
	for i := 1; i < len(b.phrases); i++ {
		resp += fmt.Sprintf(`%v CONTAINS(LOWER(s.text), LOWER("%v")) `,
			b.operations[i-1], b.phrases[i])
	}
	return resp[:len(resp)-1]
}

func NewBooleanQuery(query string) BooleanQuery {
	re := regexp.MustCompile("( AND | OR | NOT )")
	if re.MatchString(query) { //boolean operation
		operations := associateOperation(re.FindAllString(query, -1))
		phrases := re.Split(query, -1)

		return BooleanQuery{
			phrases:    phrases,
			operations: operations,
		}

	} else if sentences := strings.Split(query, ". "); len(sentences) > 1 { // multiple phrases
		operations := make([]Operation, len(sentences)-1)
		for i := 0; i < len(sentences)-1; i += 1 {
			operations[i] = AND
		}
		if last_sentence := sentences[len(sentences)-1]; string(last_sentence[len(last_sentence)-1]) == "." {
			sentences[len(sentences)-1] = last_sentence[0 : len(last_sentence)-1]
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
		if o == " AND " {
			resp[i] = AND
		} else if o == " OR " {
			resp[i] = OR
		} else {
			resp[i] = NOT
		}
	}

	return resp
}
