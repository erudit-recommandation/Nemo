package middleware

import (
	"net/http"
	"time"
)

type httpHandlerFunc func(http.ResponseWriter, *http.Request)

var LIMIT_ENTENDU_EN_VOYAGE uint = 100_000_00
var MAX_PAGE_ENTENDU_EN_VOYAGE uint = 30
var CACHE_DURATION = 1 * time.Hour

var LIMIT_RENCONTRE_ENVOYAGE uint = 20

var CACHE cache = make(cache)
