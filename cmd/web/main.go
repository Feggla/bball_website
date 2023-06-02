package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/G", guards)
	mux.HandleFunc("/C", centres)
	mux.HandleFunc("/Search?", search)
	mux.HandleFunc("/F", forwards)
	mux.HandleFunc("/users", users)
	mux.HandleFunc("/fantasyteam", fantasy)
	mux.HandleFunc("/myteam", myTeam)
	// mux.HandleFunc("/api/players", players)
	mux.HandleFunc("/api/myTeam", apiMyTeam)

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
