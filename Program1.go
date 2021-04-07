package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./assets/index.html"))
	tmpl.Execute(w, nil)
}

func main() {
	fmt.Println("Running server...")
	http.HandleFunc("/", index)
	http.ListenAndServe(":8000", nil)
}
