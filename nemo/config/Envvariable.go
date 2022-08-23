package config

import "fmt"

type EnvVariable struct {
	Port                    string           `json:"port"`
	ArangoAddr              string           `json:"arango_addr"`
	ArangoPassword          string           `json:"arango_password"`
	ArangoUsername          string           `json:"arango_username"`
	ArangoDatabase          []DatabaseCorpus `json:"arango_database"`
	TextAnalysisServiceAddr string           `json:"text_analysis_service_addr"`
}

func (e EnvVariable) IsAnExistingCorpus(corpus string) bool {
	for _, d := range e.ArangoDatabase {
		if corpus == d.Corpus {
			return true
		}
	}
	return false
}

func (e EnvVariable) GetDatabaseCorpus(corpusName string) (DatabaseCorpus, error) {
	for _, corpus := range e.ArangoDatabase {
		if corpus.Corpus == corpusName {
			return corpus, nil
		}
	}
	return DatabaseCorpus{}, fmt.Errorf("le corpus n'existe pas")
}

func (e EnvVariable) GetDatabaseCorpusExcept(corpusName string) []DatabaseCorpus {
	resp := make([]DatabaseCorpus, 0, len(e.ArangoDatabase))
	for _, corpus := range e.ArangoDatabase {
		if corpus.Name != corpusName {
			resp = append(resp, corpus)

		}

	}
	return resp
}

type DatabaseCorpus struct {
	Name   string `json:"name"`
	Corpus string `json:"corpus"`
}
