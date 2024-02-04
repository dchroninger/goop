package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

type key int

const keyServerAddr key = iota

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got / request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "Hello world!\n")
}

func getAdminRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got / request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "Hello Admin!\n")
}

func main() {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())

	// Public server
	mux := http.NewServeMux()
	port := getEnvWithDefault("PORT", "3333")

	// Public server routes
	mux.HandleFunc("/", getRoot)

	// Admin server
	adminMux := http.NewServeMux()
	adminPort := getEnvWithDefault("ADMIN_PORT", "4444")

	// Admin server routes
	adminMux.HandleFunc("/", getAdminRoot)

	publicServer := buildHttpServer(port, mux, keyServerAddr, ctx)
	adminServer := buildHttpServer(adminPort, adminMux, keyServerAddr, ctx)

	// Launch the servers in goroutines
	go listenAndServe("Public Server", publicServer, cancel)
	go listenAndServe("Admin Server", adminServer, cancel)

	<-ctx.Done()
}

func listenAndServe(name string, s *http.Server, c context.CancelFunc) {
	err := s.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("%s closed\n", name)
	} else {
		fmt.Printf("%s error: %s", name, err)
	}
	c()
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func buildHttpServer(port string, mux *http.ServeMux, ctxKey key, ctx context.Context) *http.Server {
	return &http.Server{
		Addr:    ":" + port,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			return context.WithValue(ctx, ctxKey, l.Addr().String())
		},
	}
}
