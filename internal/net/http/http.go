package http

import "net/http"

type Server = http.Server

var ErrServerClosed = http.ErrServerClosed
