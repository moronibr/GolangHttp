package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var env = "dev"
var cache map[string]*template.Template

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "index")
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "contact")
}

func RenderTemplate(w http.ResponseWriter, page string) {

	var t *template.Template
	var err error

	_, exists := cache[page]

	if !exists || env == "dev" {
		t, err = template.Must(template.ParseFiles("templates/"+page+"Template.html")).ParseFiles(
			"templates/"+page+"Template.html",
			"templates/baseTemplate.html",
		)
		if err != nil {
			log.Println(err)
			return
		}
		cache[page] = t
	} else {
		fmt.Println("Cache hit")
		t = cache[page]
	}

	t.Execute(w, nil)
}

func main() {
	cache = make(map[string]*template.Template)

	config := Config{
		Port:    "3000",
		Env:     "dev",
		Version: "1.0.0",
	}

	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/contact", ContactHandler)

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		RenderTemplate(w, "about")
	})

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil)

}
