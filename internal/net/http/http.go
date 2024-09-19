package http

import "net/http"

type (
	Dir    = http.Dir
	Server = http.Server
)

var ErrServerClosed = http.ErrServerClosed
