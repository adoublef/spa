package http

import "net/http"

func Handler() http.Handler {
	mux := http.NewServeMux()
	return mux
}
