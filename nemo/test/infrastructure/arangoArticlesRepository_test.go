package infrastructure_test

import (
	"testing"

	"github.com/erudit-recommandation/search-engine-webapp/config"
	"github.com/erudit-recommandation/search-engine-webapp/infrastructure"
)

func TestCasesArangoArticlesRepository(t *testing.T) {
	env := config.EnvVariable{
		ArangoAddr:     "http://localhost:8529",
		ArangoPassword: "rootpassword",
		ArangoUsername: "root",
		ArangoDatabase: []config.DatabaseCorpus{
			{
				Name:   "erudit",
				Corpus: "erudit",
			},
		},
	}

	config.SetConfig(&env)
	for _, ad := range env.ArangoDatabase {
		f := func() (infrastructure.ArticlesRepository, error) {
			return infrastructure.ProvideArangoArticlesRepository(ad.Corpus)
		}
		for _, testCase := range ProvideTestCaseArangoArticlesRepository() {
			test, name := testCase(f, t)
			t.Run(name, test)
		}
	}

}
