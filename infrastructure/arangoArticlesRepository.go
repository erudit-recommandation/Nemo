package infrastructure

import (
	"context"
	"fmt"
	"time"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/erudit-recommandation/search-engine-webapp/config"
	"github.com/erudit-recommandation/search-engine-webapp/domain"
)

var QUERY_MAXIMUM_DURATION = 30 * time.Second

var arangoDb *ArangoArticlesRepository = nil

type ArangoArticlesRepository struct {
	database   driver.Database
	collection string
}

func (a ArangoArticlesRepository) GetByIdproprio(id string) (domain.Article, error) {

	ctx, cancel := context.WithTimeout(context.Background(), QUERY_MAXIMUM_DURATION)
	defer cancel()

	query := fmt.Sprintf(`FOR c IN articles
    						FILTER c.idproprio == "%v"
    						LIMIT 1
    						RETURN c`,
		id)
	cursor, err := a.database.Query(ctx, query, nil)
	if err != nil {
		return domain.Article{}, err
	}
	defer cursor.Close()
	var doc domain.Article

	_, err = cursor.ReadDocument(ctx, &doc)

	if driver.IsNoMoreDocuments(err) {
		return domain.Article{}, nil
	} else if err != nil {
		return domain.Article{}, err
	}
	doc.BuildUrl()

	return doc, nil
}

func (a ArangoArticlesRepository) SearchPhrases(phrase string, n uint) ([]domain.Article, error) {

	ctx, cancel := context.WithTimeout(context.Background(), QUERY_MAXIMUM_DURATION)
	defer cancel()

	query := fmt.Sprintf(`FOR c IN article_analysis
							SEARCH ANALYSER(PHRASE(c.text,"%v"),"text_fr")
    						LIMIT %v
    						RETURN c`,
		phrase,
		n)
	cursor, err := a.database.Query(ctx, query, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	resp := make([]domain.Article, 0, n)
	for {
		var doc domain.Article
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		doc.BuildUrl()
		resp = append(resp, doc)
	}
	return resp, nil
}

func ProvideArangoArticlesRepository() (ArticlesRepository, error) {
	if arangoDb == nil {
		config := config.GetConfig()
		conn, err := http.NewConnection(http.ConnectionConfig{
			Endpoints: []string{config.ArangoPort},
			// TLSConfig: &tls.Config{ /*...*/ },
		})
		if err != nil {
			return nil, err
		}
		client, err := driver.NewClient(driver.ClientConfig{
			Connection:     conn,
			Authentication: driver.BasicAuthentication(config.ArangoUsername, config.ArangoPassword),
		})
		if err != nil {
			return nil, err
		}

		ctx := context.Background()
		db, err := client.Database(ctx, config.ArangoDatabase)
		if err != nil {
			return nil, err
		}

		return ArangoArticlesRepository{database: db,
			collection: config.ArangoArticleCollection}, nil
	}
	return arangoDb, nil

}
