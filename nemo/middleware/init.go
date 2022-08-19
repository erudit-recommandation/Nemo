package middleware

import (
	"net/http"
	"time"
)

type httpHandlerFunc func(http.ResponseWriter, *http.Request)

var LIMIT_ENTENDU_EN_VOYAGE uint = 100_000_00
var MAX_PAGE_ENTENDU_EN_VOYAGE uint = 20
var CACHE_DURATION = 1 * time.Hour

var LIMIT_RENCONTRE_EN_VOYAGE uint = 17

var CACHE_ENTENDU_EN_VOYAGE cache = make(cache)

var CACHE_RENCONTRE_EN_VOYAGE cache = make(cache)
