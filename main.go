package main

import (
	"context"
	"embed"
	"flag"
	"html/template"
	"net/http"

	"golang.org/x/net/websocket"
)

type key int

const (
	keyServerAddr key = iota
	keyHotReload  key = iota
)

var (
	templates embed.FS
	html      *template.Template
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	hotReload := flag.Bool("hotreload", false, "Enable hot reload")
	flag.Parse()

	// Hot reload websocket server
	if *hotReload {
		ctx = context.WithValue(ctx, keyHotReload, *hotReload)
		s := NewWebSocketServer()
		http.Handle("/", websocket.Handler(s.handleWS))
		go http.ListenAndServe(":5555", nil)
	}

	// Public server
	mux := http.NewServeMux()
	port := getEnvWithDefault("PORT", "3333")
	mux.HandleFunc("/css/output.css", serveCss)
	if *hotReload {
		mux.HandleFunc("/_hotreload", serveHotReload)
	}
	mux.HandleFunc("/", index)
	publicServer := buildHttpServer(port, mux, keyServerAddr, ctx)
	go listenAndServe("Public Server", publicServer, cancel)

	<-ctx.Done()
}
