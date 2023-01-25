package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/players.tmpl.html",
		"./ui/html/pages/playertable.tmpl.html",
		"./ui/html/pages/searchtable.tmpl.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	terms := r.FormValue("search")
	position := PositionFromQuery(terms)
	if position != "" {
		data, err := GetPlayersByPosition(position)
		if err != nil {
			log.Print(err)
		}
		err = ts.ExecuteTemplate(w, "searchtable", data)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
		log.Println(terms)
	} else {
		err = ts.ExecuteTemplate(w, "base", nil)

		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	}
	w.Header().Add("Content-Type", "text/plain")
}

func search(w http.ResponseWriter, r *http.Request) {
	// log.Print("Help")
	// r.ParseForm()
	// searchInput := r.Form.Get("search")
	// log.Print(searchInput)
	res, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	text := res.Query()
	search := text.Get("search")
	fmt.Println(search)

}

func guards(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/players.tmpl.html",
		"./ui/html/pages/playertable.tmpl.html",
		"./ui/html/pages/searchtable.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	data, err := Guards()
	if err != nil {
		log.Print(err)
	}
	err = ts.ExecuteTemplate(w, "players", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func centres(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/players.tmpl.html",
		"./ui/html/pages/playertable.tmpl.html",
		"./ui/html/pages/searchtable.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	data, err := Centres()
	if err != nil {
		log.Print(err)
	}
	err = ts.ExecuteTemplate(w, "players", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func forwards(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/players.tmpl.html",
		"./ui/html/pages/playertable.tmpl.html",
		"./ui/html/pages/searchtable.tmpl.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	data, err := Forwards()
	if err != nil {
		log.Print(err)
	}
	err = ts.ExecuteTemplate(w, "players", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
