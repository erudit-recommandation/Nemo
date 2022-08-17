package middleware

import "net/http"

type httpHandlerFunc func(http.ResponseWriter, *http.Request)

var LIMIT_ENTENDU_EN_VOYAGE uint = 100_00
var LIMIT_RENCONTRE_ENVOYAGE uint = 20
