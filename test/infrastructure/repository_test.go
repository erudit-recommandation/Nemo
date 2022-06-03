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
