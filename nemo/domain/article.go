package domain

type Article struct {
	Title            string `json:"title"`
	Url              string
	Year             int    `json:"annee"`
	Author           string `json:"author"`
	ID               string `json:"idproprio"`
	Journal          string `json:"titrerev"`
	CurrentSentence  string `json:"current_sentence"`
	PreviousSentence string `json:"previous_sentence"`
	NextSentence     string `json:"next_sentence"`
	RelatedText      RelatedText
	PersonaSvg       string `json:"persona_svg"`
	Bmu              int    `json:"bmu"`
}

type RelatedText struct {
	Prev  string
	Best  string
	After string
}

func (a *Article) BuildUrl(corpus string) {
	if corpus == "erudit" {
		a.Url = build_erudit_url(a)
	}
	if a.Url != "" && a.Title == "" {
		a.Title = a.Url
	}

}

func (a *Article) BuildRelatedText() {

	a.RelatedText = RelatedText{
		Prev:  a.PreviousSentence,
		Best:  a.CurrentSentence,
		After: a.NextSentence,
	}

}
