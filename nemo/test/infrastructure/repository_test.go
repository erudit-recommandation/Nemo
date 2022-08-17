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
		testSearchSentencesWithOneResult,
		testSearchSentencesWithNoResults,
		testSearchSentencesWithMultipleResults,
		testSearchSentencesWithBooleanQuery,
		testSearchSentenceWithID,
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

func testSearchSentencesWithOneResult(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	repo, err := repositoryProvider()
	if err != nil {
		t.Fatal(err)
	}

	return func(t *testing.T) {
		phrase := "l'unité et la division"
		var n uint = 1

		resp, err := repo.SearchSentences(phrase, n)
		if err != nil {
			t.Fatal(err)
		}

		if uint(len(resp)) != n {
			t.Fatalf("there should be one response")
		}

	}, "testSearchSentencesWithOneResult"
}

func testSearchSentencesWithNoResults(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	repo, err := repositoryProvider()
	if err != nil {
		t.Fatal(err)
	}

	return func(t *testing.T) {
		phrase := ""
		var n uint = 10

		resp, err := repo.SearchSentences(phrase, n)
		if err != nil {
			t.Fatal(err)
		}
		if len(resp) != 0 {
			t.Fatalf(`There should be no article with the phrase "%v" but return \n%v`, phrase, resp)
		}

	}, "testSearchSentencesWithNoResults"
}

func testSearchSentencesWithMultipleResults(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	repo, err := repositoryProvider()
	if err != nil {
		t.Fatal(err)
	}

	return func(t *testing.T) {
		phrase := "logements"
		var n uint = 20

		resp, err := repo.SearchSentences(phrase, n)
		if err != nil {
			t.Fatal(err)
		}

		if uint(len(resp)) > n || uint(len(resp)) == 0 {
			t.Fatalf("there should be between 0 and %v results", n)
		}

	}, "testSearchSentencesWithMultipleResults"
}

func testSearchSentencesWithBooleanQuery(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	repo, err := repositoryProvider()
	if err != nil {
		t.Fatal(err)
	}

	return func(t *testing.T) {
		query := "Le voyage AND La ville de quebec OR le retour NOT l'economie du congo"
		var n uint = 20

		resp, err := repo.SearchSentences(query, n)
		if err != nil {
			t.Fatal(err)
		}

		if uint(len(resp)) > n || uint(len(resp)) == 0 {
			t.Fatalf("there should be between 0 and %v results", n)
		}

	}, "testSearchSentencesWithBooleanQuery"
}

func testSearchSentenceWithID(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	repo, err := repositoryProvider()
	if err != nil {
		t.Fatal(err)
	}

	return func(t *testing.T) {
		query := "Le voyage AND La ville de quebec OR le retour NOT l'economie du congo"
		var n uint = 20

		resp, err := repo.GetSearchSentencesID(query, n)
		if err != nil {
			t.Fatal(err)
		}

		if uint(len(resp)) > n || uint(len(resp)) == 0 {
			t.Fatalf("there should be between 0 and %v results", n)
		}

		for _, id := range resp {
			_, err = repo.GetArticleFromSentenceID(id)

			if err != nil {
				t.Fatal(err)
			}
		}

	}, "testSearchSentenceWithID"
}
