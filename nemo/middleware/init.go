package middleware

import "net/http"

type httpHandlerFunc func(http.ResponseWriter, *http.Request)

var LIMIT uint = 30
