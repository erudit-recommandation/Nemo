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

func (a ArangoArticlesRepository) GetByTitle(title string, n uint) (domain.Article, error) {

	ctx, cancel := context.WithTimeout(context.Background(), QUERY_MAXIMUM_DURATION)
	defer cancel()

	query := fmt.Sprintf(`FOR c IN articles
    						FILTER c.title == "%v"
    						LIMIT %v
    						RETURN c`,
		title,
		n)
	cursor, err := a.database.Query(ctx, query, nil)
	if err != nil {
		return domain.Article{}, err
	}
	defer cursor.Close()
	var doc domain.Article
	for {
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return domain.Article{}, err
		}
	}
	return doc, nil
}

func (a ArangoArticlesRepository) GetByAuthor(title string, n uint) ([]domain.Article, error) {

	ctx, cancel := context.WithTimeout(context.Background(), QUERY_MAXIMUM_DURATION)
	defer cancel()

	query := fmt.Sprintf(`FOR c IN articles
    						FILTER c.author == "%v"
    						LIMIT %v
    						RETURN c`,
		title,
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
