package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
)

const keyServerAddr = "serverAddr"

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got / request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "Hello world!\n")
}

func GetHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "Hello, HTTP!\n")
}

func main() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", GetHello)

	ctx, cancel := context.WithCancel(context.Background())
	publicServer := &http.Server{
		Addr:    ":3333",
		Handler: http.DefaultServeMux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	adminServer := &http.Server{
		Addr:    ":4444",
		Handler: http.DefaultServeMux,
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
	go func() {
		err := adminServer.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("admin server closed")
		} else {
			fmt.Println("admin server error:", err)
		}
		cancel()
	}()
}
