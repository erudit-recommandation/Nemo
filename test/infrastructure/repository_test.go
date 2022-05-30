package infrastructure_test

import (
	"testing"

	"github.com/erudit-recommandation/search-engine-webapp/domain"
	"github.com/erudit-recommandation/search-engine-webapp/infrastructure"
	"github.com/erudit-recommandation/search-engine-webapp/test"
)

func ProvideTestCaseArangoArticlesRepository() []func(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	return []func(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string){
		testGetAuthorArticles,
		testGetAuthorArticlesAuthorDontExist,
		testGetTitles,
		testGetTitleArticlesTitleDontExist,
	}
}

func testGetAuthorArticles(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	repo, err := repositoryProvider()
	if err != nil {
		t.Fatal(err)
	}

	return func(t *testing.T) {
		var n uint = 10
		author := "Savoie-Bernard"
		expectedResult := []domain.Article{
			{
				Title: "Hervé : dans l’antichambre du véritable amour",
				Year:  2021,
				ID:    "96127ac",
			},
			{
				Year:  2017,
				ID:    "86259ac",
				Title: "Des bouchées de Chloé",
			},
			{
				Year:  2021,
				ID:    "96215ac",
				Title: "Réflexion.  Qui porte les voix dans l’écriture de l’histoire ?",
			},
			{
				Year:  2016,
				Title: "Se réinventer with a new tab",
				ID:    "1044407ar",
			},
			{
				Year:  2021,
				Title: "Lettre à un·e écrivain·e vivant·e malgré tout",
				ID:    "97262ac",
			},
		}

		for _, el := range expectedResult {
			el.Author = author
			el.BuildUrl()
		}

		resp, err := repo.GetByAuthor(author, n)
		if err != nil {
			t.Fatal(err)
		}
		if len(resp) != len(expectedResult) {
			t.Fatalf("should return %v articles but returned %v", len(expectedResult), resp)
		}
		if !test.SliceContainTheSameElements(resp, expectedResult) {
			t.Fatalf("the articles are not same the response is %v\n where we expect \n %v", resp, expectedResult)
		}

	}, "testGetAuthorArticles"
}

func testGetAuthorArticlesAuthorDontExist(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	repo, err := repositoryProvider()
	if err != nil {
		t.Fatal(err)
	}

	return func(t *testing.T) {
		var n uint = 10
		author := "Carlos Marximilian"

		resp, err := repo.GetByAuthor(author, n)
		if err != nil {
			t.Fatal(err)
		}
		if len(resp) != 0 {
			t.Fatalf("There should be no author of the name %v", author)
		}

	}, "testGetAuthorArticlesAuthorDontExist"
}

func testGetTitles(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	repo, err := repositoryProvider()
	if err != nil {
		t.Fatal(err)
	}

	return func(t *testing.T) {
		var n uint = 10
		title := "L’injonction de penser dans l’amitié"
		expectedResult := []domain.Article{
			{
				Title: "L’injonction de penser dans l’amitié",
				Year:  2005,
				ID:    "18424ac",
			},
		}

		for _, el := range expectedResult {
			el.BuildUrl()
		}

		resp, err := repo.GetByTitle(title, n)
		if err != nil {
			t.Fatal(err)
		}
		if len(resp) != len(expectedResult) {
			t.Fatalf("should return %v articles but returned %v", len(expectedResult), resp)
		}
		if !test.SliceContainTheSameElements(resp, expectedResult) {
			t.Fatalf("the articles are not same the response is %v\n where we expect \n %v", resp, expectedResult)
		}

	}, "testGetTitles"
}

func testGetTitleArticlesTitleDontExist(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	repo, err := repositoryProvider()
	if err != nil {
		t.Fatal(err)
	}

	return func(t *testing.T) {
		var n uint = 10
		title := "Switch and translation"

		resp, err := repo.GetByTitle(title, n)
		if err != nil {
			t.Fatal(err)
		}
		if len(resp) != 0 {
			t.Fatalf("There should be no article with the title %v", title)
		}

	}, "testGetTitleArticlesTitleDontExist"
}
