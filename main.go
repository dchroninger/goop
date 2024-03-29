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
	keyWatch      key = iota
	keyPort       key = iota
)

var (
	templates embed.FS
	html      *template.Template
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	watch := flag.Bool("watch", false, "Enable watch for browser refresh")
	port := flag.Int("port", 3333, "Port for the HTTP server to run on")
	flag.Parse()

	// Hot reload websocket server
	if *watch {
		ctx = context.WithValue(ctx, keyWatch, *watch)
		// Port is only needed for live reload
		ctx = context.WithValue(ctx, keyPort, *port)
		s := NewWebSocketServer()
		http.Handle("/", websocket.Handler(s.handleWS))
		go http.ListenAndServe(":5555", nil)
	}

	// Public server
	mux := http.NewServeMux()
	mux.HandleFunc("/css/output.css", serveCss)
	if *watch {
		mux.HandleFunc("/_hotreload", serveHotReload)
	}
	mux.HandleFunc("/", index)
	publicServer := buildHttpServer(*port, mux, keyServerAddr, ctx)
	go listenAndServe("Public Server", publicServer, cancel)

	<-ctx.Done()
}
