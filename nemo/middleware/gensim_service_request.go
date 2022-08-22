package middleware

type gemsimServiceRequest struct {
	Text   string `json:"text"`
	N      uint   `json:"n"`
	Corpus string `json:"corpus"`
}
