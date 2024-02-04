package main

import (
	"context"
	"embed"
	"html/template"
	"net/http"
)

type key int

const keyServerAddr key = iota

var (
	templates embed.FS
	css       embed.FS
	html      *template.Template
)

func main() {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())

	// Public server
	mux := http.NewServeMux()
	port := getEnvWithDefault("PORT", "3333")
	mux.Handle("/css/output.css", http.FileServer(http.FS(css)))
	mux.HandleFunc("/", index)
	publicServer := buildHttpServer(port, mux, keyServerAddr, ctx)
	go listenAndServe("Public Server", publicServer, cancel)

	<-ctx.Done()
}
