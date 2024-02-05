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
	HotReload bool
}

func index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	hotReload, _ := ctx.Value(keyHotReload).(bool)
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, IndexData{HotReload: hotReload})
}
