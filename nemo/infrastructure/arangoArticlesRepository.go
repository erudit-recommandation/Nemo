package infrastructure

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"sort"
	"time"

	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/erudit-recommandation/search-engine-webapp/config"
	"github.com/erudit-recommandation/search-engine-webapp/domain"
)

var arangoDb map[string]*ArangoArticlesRepository = make(map[string]*ArangoArticlesRepository)

type ArangoArticlesRepository struct {
	database driver.Database
}

func (a ArangoArticlesRepository) GetByIdproprio(id string) (domain.Article, error) {

	ctx, cancel := context.WithTimeout(context.Background(), QUERY_MAXIMUM_DURATION)
	defer cancel()

	query := fmt.Sprintf(`FOR a IN articles
    						FILTER a.idproprio == "%v"
    						LIMIT 1
    						RETURN {
								title:a.title,
								annee:a.annee,
								author:a.author,
								idproprio:a.idproprio,
								titrerev:a.titrerev,
								persona_svg: a.persona_svg,
								bmu:a.bmu
						}`,
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

func (a ArangoArticlesRepository) GetByIdPandas(id int) (domain.Article, error) {

	ctx, cancel := context.WithTimeout(context.Background(), QUERY_MAXIMUM_DURATION)
	defer cancel()

	query := fmt.Sprintf(`FOR a IN articles
    						FILTER a.pandas_index == %v
    						LIMIT 1
    						RETURN {
								title:a.title,
								annee:a.annee,
								author:a.author,
								idproprio:a.idproprio,
								titrerev:a.titrerev,
								persona_svg: a.persona_svg
						}`,
		id)
	cursor, err := a.database.Query(ctx, query, nil)
	if err != nil {
		return domain.Article{}, err
	}
	defer cursor.Close()
	var doc domain.Article

	_, err = cursor.ReadDocument(ctx, &doc)
	if reflect.DeepEqual(doc, domain.Article{}) {
		return domain.Article{}, fmt.Errorf("was not able to find an article")
	}

	if driver.IsNoMoreDocuments(err) {
		return domain.Article{}, nil
	} else if err != nil {
		return domain.Article{}, err
	}
	doc.BuildUrl()

	return doc, nil
}

func (a ArangoArticlesRepository) SearchSentences(phrase string, n uint) ([]domain.Article, error) {

	ctx, cancel := context.WithTimeout(context.Background(), QUERY_MAXIMUM_DURATION)
	defer cancel()
	booleanQuery := domain.NewBooleanQuery(phrase)

	query := fmt.Sprintf(`LET ids_sentences = (
FOR s IN article_analysis
SEARCH ANALYZER(
	%v,
	"text_fr"
)
    LIMIT %v
    RETURN {"id":s.idproprio, n_sentence:s.index_nm, sentence: s.text}
    
)


FOR el IN ids_sentences

    LET prev = (
        FOR s IN sentences
            FILTER el.n_sentence == s.index_nm -1 AND el.id == s.idproprio
            LIMIT 1
            RETURN s.text
    
    )
    
    LET next = (
        FOR s IN sentences
            FILTER el.n_sentence == s.index_nm +1 AND el.id == s.idproprio
            LIMIT 1
            RETURN s.text
    
    )

    FOR a IN articles
        FILTER a.idproprio == el.id
        RETURN DISTINCT {title:a.title,
                annee:a.annee,
                author:a.author,
                idproprio:a.idproprio,
                titrerev:a.titrerev,
                current_sentence:el.sentence,
                previous_sentence: prev[0],
                next_sentence: next[0],
				persona_svg: a.persona_svg
        }`,
		booleanQuery.ToArangoPhraseQueryBody(),
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
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(resp), func(i, j int) { resp[i], resp[j] = resp[j], resp[i] })
	return resp, nil
}

func (a ArangoArticlesRepository) GetSearchSentencesID(phrase string, n uint) ([]ArticlesID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), QUERY_MAXIMUM_DURATION)
	defer cancel()
	booleanQuery := domain.NewBooleanQuery(phrase)

	query := fmt.Sprintf(`
FOR s IN article_analysis
SEARCH ANALYZER(
	%v,
	"text_fr"
)
    LIMIT %v
    RETURN {"id":s.idproprio, n_sentence:s.index_nm}`,
		booleanQuery.ToArangoPhraseQueryBody(),
		n)
	cursor, err := a.database.Query(ctx, query, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	resp := make([]ArticlesID, 0, n)
	for {
		var doc ArticlesID
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, err
		}
		resp = append(resp, doc)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(resp), func(i, j int) { resp[i], resp[j] = resp[j], resp[i] })
	return resp, nil
}

func (a ArangoArticlesRepository) GetArticleFromSentenceID(articleID ArticlesID) (domain.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), QUERY_MAXIMUM_DURATION)
	defer cancel()

	query := fmt.Sprintf(`
	LET index = %v
	LET id = "%v"	
	LET prev = (
        FOR s IN sentences
            FILTER index == s.index_nm -1  AND id == s.idproprio
            LIMIT 1
            RETURN s.text
    
    )
    
    LET next = (
        FOR s IN sentences
            FILTER index == s.index_nm +1 AND id == s.idproprio
            LIMIT 1
            RETURN s.text
    
    )
    
     LET current = (
        FOR s IN sentences
            FILTER index == s.index_nm AND id == s.idproprio
            LIMIT 1
            RETURN s.text
    
    )

    FOR a IN articles
        FILTER a.idproprio == id
        LIMIT 1
        RETURN {title:a.title,
                annee:a.annee,
                author:a.author,
                idproprio:a.idproprio,
                titrerev:a.titrerev,
                current_sentence:current[0],
                previous_sentence: prev[0],
                next_sentence: next[0],
				persona_svg: a.persona_svg,
				bmu:a.bmu
        }
	`,
		articleID.NSentence, articleID.ID)
	cursor, err := a.database.Query(ctx, query, nil)
	if err != nil {
		return domain.Article{}, err
	}
	defer cursor.Close()
	var doc domain.Article

	_, err = cursor.ReadDocument(ctx, &doc)
	if reflect.DeepEqual(doc, domain.Article{}) {
		return domain.Article{}, fmt.Errorf("was not able to find an article %v", articleID)
	}

	if driver.IsNoMoreDocuments(err) {
		return domain.Article{}, nil
	} else if err != nil {
		return domain.Article{}, err
	}
	doc.BuildUrl()

	return doc, nil
}

func (a ArangoArticlesRepository) GetNeighbouringArticlesByBMU(bmu int, limit uint) ([]domain.Article, error) {
	ctx, cancel := context.WithTimeout(context.Background(), QUERY_MAXIMUM_DURATION)
	defer cancel()
	query := fmt.Sprintf(`
	LET bmu = %v
	FOR a IN articles
		FILTER a.bmu==bmu OR a.bmu==bmu+1 OR a.bmu==bmu-1
		LIMIT %v
		RETURN {
			title:a.title,
			annee:a.annee,
			author:a.author,
			idproprio:a.idproprio,
			titrerev:a.titrerev,
			persona_svg: a.persona_svg,
			bmu:a.bmu
		}
		
	`, bmu, limit)

	cursor, err := a.database.Query(ctx, query, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	resp := make([]domain.Article, 0, limit)
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
	sort.SliceStable(resp, func(i, j int) bool {
		return math.Abs(float64(resp[i].Bmu-bmu)) < math.Abs(float64(resp[j].Bmu-bmu))
	})
	return resp, nil
}

func ProvideArangoArticlesRepository(corpus string) (ArticlesRepository, error) {
	if !config.GetConfig().IsAnExistingCorpus(corpus) {
		return nil, fmt.Errorf("the corpus \"%v\"does not exist", corpus)
	}
	if _, ok := arangoDb[corpus]; !ok {
		config := config.GetConfig()
		conn, err := http.NewConnection(http.ConnectionConfig{
			Endpoints: []string{config.ArangoAddr},
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
		db, err := client.Database(ctx, corpus)
		if err != nil {
			return nil, err
		}

		return ArangoArticlesRepository{database: db}, nil
	}
	return arangoDb[corpus], nil

}
