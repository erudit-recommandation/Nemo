package infrastructure_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/erudit-recommandation/search-engine-webapp/domain"
	"github.com/erudit-recommandation/search-engine-webapp/infrastructure"
	"github.com/erudit-recommandation/search-engine-webapp/test"
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
			Title:  "Imperturbablement... Au sujet d’une nouvelle génération intellectuelle",
			Year:   2005,
			Author: "Thibault",
			ID:     id,
		}
		expectedResult.BuildUrl()

		resp, err := repo.GetByIdproprio(id)
		if err != nil {
			t.Fatal(err)
		}

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
		var n uint = 10

		expectedResp := []domain.Article{
			{
				Author: "Hamel",
				Title:  "La question du partenariat : de la crise institutionnelle à la redéfinition des rapports entre sphère publique et sphère privée",
				Year:   1995,
				ID:     "1002278ar",
			},
		}
		expectedResp[0].BuildUrl()

		resp, err := repo.SearchPhrases(phrase, n)
		if err != nil {
			t.Fatal(err)
		}

		if !test.SliceContainTheSameElements(resp, expectedResp) {
			t.Fatalf("the response are not the same the response is \n%v\n whereas the expected result is\n%v",
				resp, expectedResp)
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
		var n uint = 3

		expectedResp := []domain.Article{
			{
				Author: "Altenor",
				Title:  "Les causes économiques et socio-politiques du passage de la régionalisation à la départementalisation en Haïti",
				Year:   2020,
				ID:     "1075860ar",
			},
			{
				Author: "Gherghel",
				Title:  "Transformations de la régulation politique et juridique de la famille. La Roumanie dans la période communiste et post-communiste",
				Year:   2006,
				ID:     "015786ar",
			},
			{
				Author: "Semmoud",
				Title:  "Appropriations et usages des espaces urbains en Algérie du Nord",
				Year:   2009,
				ID:     "038144ar",
			},
		}
		for i := 0; i < len(expectedResp); i++ {
			expectedResp[i].BuildUrl()
		}

		fmt.Println(expectedResp[0])
		fmt.Println("\n\na")

		resp, err := repo.SearchPhrases(phrase, n)
		if err != nil {
			t.Fatal(err)
		}

		if !test.SliceContainTheSameElements(resp, expectedResp) {
			t.Fatalf("the response are not the same the response is \n%v\n whereas the expected result is\n%v",
				resp, expectedResp)
		}

	}, "testSearchPhrasesWithMultipleResults"
}
