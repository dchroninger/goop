package main

import (
	"html/template"
	"net/http"
)

func css(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "css/output.css")
}

func index(w http.ResponseWriter, _ *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}
