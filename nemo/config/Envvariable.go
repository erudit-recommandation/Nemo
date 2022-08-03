package config

type EnvVariable struct {
	Port                    string
	ArangoPort              string
	ArangoPassword          string
	ArangoUsername          string
	ArangoDatabase          string
	ArangoArticleCollection string
	TEXT_ANALYSIS_SERVICE   string
}
