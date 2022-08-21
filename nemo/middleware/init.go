package middleware

import (
	"net/http"
	"time"
)

type httpHandlerFunc func(http.ResponseWriter, *http.Request)

var LIMIT_ENTENDU_EN_VOYAGE uint = 100_000
var MAX_BY_PAGE_ENTENDU_EN_VOYAGE uint = 20
var CACHE_DURATION = 1 * time.Hour

var LIMIT_RENCONTRE_EN_VOYAGE uint = 16
var LIMIT_ACCOSTER_EN_VOYAGE uint = 1000

var CACHE_ENTENDU_EN_VOYAGE cache = make(cache)

var CACHE_RENCONTRE_EN_VOYAGE cache = make(cache)

var CACHE_ACCOSTE_EN_VOYAGE cache = make(cache)
