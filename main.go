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

const (
	keyServerAddr key = iota
)

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
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/", getAdminRoot)

	ctx, cancel := context.WithCancel(context.Background())

	port := getEnvWithDefault("PORT", "3333")
	adminPort := getEnvWithDefault("ADMIN_PORT", "4444")

	publicServer := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	go func() {
		err := publicServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("public server closed")
		} else {
			fmt.Println("public server error:", err)
		}
		cancel()
	}()

	adminServer := &http.Server{
		Addr:    ":" + adminPort,
		Handler: adminMux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	go func() {
		err := adminServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("admin server closed")
		} else {
			fmt.Println("admin server error:", err)
		}
		cancel()
	}()

	<-ctx.Done()
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
