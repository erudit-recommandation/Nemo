package infrastructure_test

import (
	"testing"

	"github.com/erudit-recommandation/search-engine-webapp/config"
	"github.com/erudit-recommandation/search-engine-webapp/infrastructure"
)

func TestCasesArangoArticlesRepository(t *testing.T) {
	env := config.EnvVariable{
		ArangoPort:              "http://localhost:8529",
		ArangoPassword:          "rootpassword",
		ArangoUsername:          "root",
		ArangoDatabase:          "erudit",
		ArangoArticleCollection: "articles",
	}

	config.SetConfig(&env)
	for _, testCase := range ProvideTestCaseArangoArticlesRepository() {
		test, name := testCase(infrastructure.ProvideArangoArticlesRepository, t)
		t.Run(name, test)
	}
}
