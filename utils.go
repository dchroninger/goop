package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

func listenAndServe(name string, s *http.Server, c context.CancelFunc) {
	err := s.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("%s closed\n", name)
	} else {
		fmt.Printf("%s error: %s", name, err)
	}
	c()
}

func buildHttpServer(port int, mux *http.ServeMux, ctxKey key, ctx context.Context) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			return context.WithValue(ctx, ctxKey, l.Addr().String())
		},
	}
}
