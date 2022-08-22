package config

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

func (e EnvVariable) GetCorpusNames() []string {
	var names []string = make([]string, len(e.ArangoDatabase))
	for i, corpus := range e.ArangoDatabase {
		names[i] = corpus.Name
	}
	return names
}

type DatabaseCorpus struct {
	Name   string `json:"name"`
	Corpus string `json:"corpus"`
}
