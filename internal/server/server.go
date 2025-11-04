package server

import (
	"net/http"
	"time"
)

func NewServer(mux *http.ServeMux) *http.Server {
	return &http.Server{
		Addr:         ":80",
		IdleTimeout:  time.Duration(30 * time.Second),
		ReadTimeout:  time.Duration(30 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),

		Handler: mux,
	}
}
