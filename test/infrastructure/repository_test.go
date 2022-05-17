package infrastructure_test

import (
	"testing"

	"github.com/erudit-recommandation/search-engine-webapp/infrastructure"
)

func ProvideTestCaseArangoArticlesRepository() []func(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string) {
	return []func(repositoryProvider func() (infrastructure.ArticlesRepository, error), t *testing.T) (func(t *testing.T), string){
		testGetAuthorArticles,
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
		nArticle := 2
		resp, err := repo.GetByAuthor(author, n)
		if err != nil {
			t.Fatal(err)
		}
		if len(resp) != nArticle {
			t.Fatalf("should return 2 articles but returned %v", resp)
		}

	}, "testGetAuthorArticles"
}
