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
	html      *template.Template
)

func main() {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())

	// Public server
	mux := http.NewServeMux()
	port := getEnvWithDefault("PORT", "3333")
	mux.HandleFunc("/css/output.css", css)
	mux.HandleFunc("/", index)
	publicServer := buildHttpServer(port, mux, keyServerAddr, ctx)
	go listenAndServe("Public Server", publicServer, cancel)

	<-ctx.Done()
}
