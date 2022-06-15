package middleware

type gemsimServiceRequest struct {
	Text string `json:"text"`
	N    int    `json:"n"`
}
