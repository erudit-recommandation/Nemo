package infrastructure_test

import (
	"reflect"
	"testing"

	"github.com/erudit-recommandation/search-engine-webapp/domain"
	"github.com/erudit-recommandation/search-engine-webapp/infrastructure"
)

func ProvideTestCaseArangoArticlesRepository() []func(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	return []func(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string){
		testGetByIdproprio,
		testGetByIdproprioArticleDontExist,
		testSearchPhrasesWithOneResult,
		testSearchPhrasesWithNoResults,
		testSearchPhrasesWithMultipleResults,
	}
}

func testGetByIdproprio(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	repo, err := repositoryProvider()
	if err != nil {
		t.Fatal(err)
	}

	return func(t *testing.T) {
		id := "18411ac"
		expectedResult := domain.Article{
			Title:   "Imperturbablement... Au sujet d’une nouvelle génération intellectuelle",
			Year:    2005,
			Author:  "Thibault",
			ID:      id,
			Journal: "Spirale",
		}
		expectedResult.BuildUrl()

		resp, err := repo.GetByIdproprio(id)
		if err != nil {
			t.Fatal(err)
		}

		resp.Text = ""

		if !reflect.DeepEqual(resp, expectedResult) {
			t.Fatalf("the articles are not same the response is \n%v\n where we expect \n %v", resp, expectedResult)
		}

	}, "testGetByIdproprio"
}

func testGetByIdproprioArticleDontExist(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	repo, err := repositoryProvider()
	if err != nil {
		t.Fatal(err)
	}

	return func(t *testing.T) {
		id := "abcdefghi"

		resp, err := repo.GetByIdproprio(id)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(resp, domain.Article{}) {
			t.Fatalf("There should be no article with the id %v", id)
		}

	}, "testGetByIdproprioArticleDontExist"
}

func testSearchPhrasesWithOneResult(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	repo, err := repositoryProvider()
	if err != nil {
		t.Fatal(err)
	}

	return func(t *testing.T) {
		phrase := "L'unité et la division"
		var n uint = 1

		resp, err := repo.SearchPhrases(phrase, n)
		if err != nil {
			t.Fatal(err)
		}

		if uint(len(resp)) != n {
			t.Fatalf("there should be one response")
		}

	}, "testSearchPhrasesWithOneResult"
}

func testSearchPhrasesWithNoResults(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	repo, err := repositoryProvider()
	if err != nil {
		t.Fatal(err)
	}

	return func(t *testing.T) {
		phrase := ""
		var n uint = 10

		resp, err := repo.SearchPhrases(phrase, n)
		if err != nil {
			t.Fatal(err)
		}
		if len(resp) != 0 {
			t.Fatalf(`There should be no article with the phrase "%v" but return \n%v`, phrase, resp)
		}

	}, "testSearchPhrasesWithNoResults"
}

func testSearchPhrasesWithMultipleResults(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	repo, err := repositoryProvider()
	if err != nil {
		t.Fatal(err)
	}

	return func(t *testing.T) {
		phrase := "la crise des logements"
		var n uint = 20

		resp, err := repo.SearchPhrases(phrase, n)
		if err != nil {
			t.Fatal(err)
		}

		if uint(len(resp)) > n || uint(len(resp)) == 0 {
			t.Fatalf("there should be between 0 and %v results", n)
		}

	}, "testSearchPhrasesWithMultipleResults"
}
