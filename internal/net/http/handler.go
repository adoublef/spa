package http

import (
	"net/http"
)

func Handler(dir Dir) http.Handler {
	mux := http.NewServeMux()
	handleFunc := func(pattern string, h http.Handler) {
		mux.Handle(pattern, h)
	}
	handleFunc("/", http.FileServer(dir))
	return mux
}
