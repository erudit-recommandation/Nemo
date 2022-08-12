package domain_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/erudit-recommandation/search-engine-webapp/domain"
)

func TestQueryOnePhraseWithNoOperator(t *testing.T) {
	query := "le vol d'un oiseau"
	booleanQuery := domain.NewBooleanQuery(query)
	respPhrases := []string{query}
	respArangoQuery := fmt.Sprintf(`CONTAINS(LOWER(s.text), LOWER("%v"))`, query)

	if phrase := booleanQuery.Phrase(); !reflect.DeepEqual(phrase, respPhrases) {
		t.Errorf("boolean query has the phrase\n%v\n but suppose to have\n%v",
			phrase, respPhrases)
	}

	if len(booleanQuery.Operations()) != 0 {
		t.Error("should not have any operator")
	}

	if arangoQuery := booleanQuery.ToArangoPhraseQueryBody(); arangoQuery != respArangoQuery {
		t.Errorf("expected\n %v but received\n %v", respArangoQuery, arangoQuery)
	}

}

func TestQueryMultiplePhrase(t *testing.T) {
	query := "le vol d'un oiseau. Le chant d'une cygogne. 1.1+1.2=2.3."
	booleanQuery := domain.NewBooleanQuery(query)
	respPhrases := []string{"le vol d'un oiseau",
		"Le chant d'une cygogne",
		"1.1+1.2=2.3"}
	respOperations := []domain.Operation{domain.AND, domain.AND}
	respArangoQuery := fmt.Sprintf(`CONTAINS(LOWER(s.text), LOWER("%v")) AND CONTAINS(LOWER(s.text), LOWER("%v")) AND CONTAINS(LOWER(s.text), LOWER("%v"))`,
		respPhrases[0], respPhrases[1], respPhrases[2])

	if phrase := booleanQuery.Phrase(); !reflect.DeepEqual(phrase, respPhrases) {
		t.Errorf("boolean query has the phrases\n%#v\n but suppose to have\n%#v",
			phrase, respPhrases)
	}

	if operations := booleanQuery.Operations(); !reflect.DeepEqual(operations, respOperations) {
		t.Errorf("boolean query has the operations\n%#v\n but suppose to have\n%#v",
			operations, respOperations)
	}

	if arangoQuery := booleanQuery.ToArangoPhraseQueryBody(); arangoQuery != respArangoQuery {
		t.Errorf("expected\n %v but received\n %v", respArangoQuery, arangoQuery)
	}
}

func TestQueryWithSimpleOperation(t *testing.T) {

	var tests = []struct {
		query           string
		respPhrases     []string
		respOperations  []domain.Operation
		respArangoQuery string
	}{
		{"le vol d'un oiseau AND Le chant d'une cygogne",
			[]string{"le vol d'un oiseau",
				"Le chant d'une cygogne",
			},
			[]domain.Operation{domain.AND},
			fmt.Sprintf(`CONTAINS(LOWER(s.text), LOWER("%v")) AND CONTAINS(LOWER(s.text), LOWER("%v"))`, "le vol d'un oiseau", "Le chant d'une cygogne"),
		},
		{"Le chat et le renard OR Le chien et le loup",
			[]string{"Le chat et le renard",
				"Le chien et le loup",
			},
			[]domain.Operation{domain.OR},
			fmt.Sprintf(`CONTAINS(LOWER(s.text), LOWER("%v")) OR CONTAINS(LOWER(s.text), LOWER("%v"))`, "Le chat et le renard", "Le chien et le loup"),
		},

		{"abcde NOT bar",
			[]string{"abcde",
				"bar",
			},
			[]domain.Operation{domain.NOT},
			fmt.Sprintf(`CONTAINS(LOWER(s.text), LOWER("%v")) AND NOT CONTAINS(LOWER(s.text), LOWER("%v"))`, "abcde", "bar"),
		},
	}
	for _, tt := range tests {
		booleanQuery := domain.NewBooleanQuery(tt.query)

		if phrase := booleanQuery.Phrase(); !reflect.DeepEqual(phrase, tt.respPhrases) {
			t.Errorf("boolean query has the phrases\n%#v\n but suppose to have\n%#v",
				phrase, tt.respPhrases)
		}

		if operations := booleanQuery.Operations(); !reflect.DeepEqual(operations, tt.respOperations) {
			t.Errorf("boolean query has the operations\n%#v\n but suppose to have\n%#v",
				operations, tt.respOperations)
		}

		if arangoQuery := booleanQuery.ToArangoPhraseQueryBody(); arangoQuery != tt.respArangoQuery {
			t.Errorf("expected\n %v but received\n %v", tt.respArangoQuery, arangoQuery)
		}
	}
}

func TestQueryWithCombinedOperation(t *testing.T) {

	var tests = []struct {
		query           string
		respPhrases     []string
		respOperations  []domain.Operation
		respArangoQuery string
	}{
		{"le vol d'un oiseau AND Le chant d'une cygogne NOT FOO",
			[]string{"le vol d'un oiseau",
				"Le chant d'une cygogne",
				"FOO",
			},
			[]domain.Operation{domain.AND, domain.NOT},
			fmt.Sprintf(`CONTAINS(LOWER(s.text), LOWER("%v")) AND CONTAINS(LOWER(s.text), LOWER("%v")) AND NOT CONTAINS(LOWER(s.text), LOWER("%v"))`, "le vol d'un oiseau", "Le chant d'une cygogne", "FOO"),
		},
		{"Le chat et le renard OR Le chien et le loup AND a OR b NOT c",
			[]string{"Le chat et le renard",
				"Le chien et le loup",
				"a",
				"b",
				"c",
			},
			[]domain.Operation{domain.OR, domain.AND, domain.OR, domain.NOT},
			fmt.Sprintf(`CONTAINS(LOWER(s.text), LOWER("%v")) OR CONTAINS(LOWER(s.text), LOWER("%v")) AND CONTAINS(LOWER(s.text), LOWER("%v")) OR CONTAINS(LOWER(s.text), LOWER("%v")) AND NOT CONTAINS(LOWER(s.text), LOWER("%v"))`,
				"Le chat et le renard", "Le chien et le loup", "a", "b", "c"),
		},
	}
	for _, tt := range tests {
		booleanQuery := domain.NewBooleanQuery(tt.query)

		if phrase := booleanQuery.Phrase(); !reflect.DeepEqual(phrase, tt.respPhrases) {
			t.Errorf("boolean query has the phrases\n%#v\n but suppose to have\n%#v",
				phrase, tt.respPhrases)
		}

		if operations := booleanQuery.Operations(); !reflect.DeepEqual(operations, tt.respOperations) {
			t.Errorf("boolean query has the operations\n%#v\n but suppose to have\n%#v",
				operations, tt.respOperations)
		}
	}
}
