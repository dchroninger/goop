package main

import (
	"html/template"
	"net/http"
)

func serveHotReload(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./hotreload.js")
}

func serveCss(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "css/output.css")
}

type IndexData struct {
	Watch bool
	Port  int
}

func index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	watch, _ := ctx.Value(keyWatch).(bool)
	port, _ := ctx.Value(keyPort).(int)

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, IndexData{Watch: watch, Port: port})
}
