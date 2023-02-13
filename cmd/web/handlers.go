package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	_ "github.com/lib/pq"
)

var files = []string{
	"./ui/html/base.tmpl.html",
	"./ui/html/pages/home.tmpl.html",
	"./ui/html/partials/nav.tmpl.html",
	"./ui/html/pages/players.tmpl.html",
	"./ui/html/pages/playertable.tmpl.html",
	"./ui/html/pages/searchtable.tmpl.html",
	"./ui/html/pages/users.tmpl.html",
	"./ui/html/pages/usertable.tmpl.html",
	"./ui/html/pages/fantasy.tmpl.html",
	"./ui/html/pages/fantasytable.tmpl.html",
	"./ui/html/pages/myteam.tmpl.html",
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	terms := r.FormValue("search")
	position := PosFromQuery(terms)
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

func users(w http.ResponseWriter, r *http.Request) {
	const (
		host     = "localhost"
		user     = "newuser"
		password = "password"
		dbname   = "postgres"
	)

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbname, password, host)
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err != nil {
		log.Print(err)
		panic(err)
	}
	fmt.Printf("\nSuccessfully connected to database!\n")

	var (
		userID int64
		name   string
	)
	type userind struct {
		Id       int64
		Username string
	}
	rows, err := db.Query("SELECT id, username FROM users order by id asc")
	if err != nil {
		panic(err)
	}
	var users []userind
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&userID, &name)
		if err != nil {
			panic(err)
		}
		users = append(users, userind{
			Id:       userID,
			Username: name,
		})

		fmt.Println(users)
		fmt.Println("\n", userID, name)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.ExecuteTemplate(w, "usertable", users)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}
func fantasy(w http.ResponseWriter, r *http.Request) {
	res, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	text := res.Query()
	username := text.Get("user")
	id := res.Query()
	addId := id.Get("ID")
	fmt.Println(username)
	switch username {
	case "":
		fmt.Println("Empty string")
		data, err := AllFantasyPlayers()
		if err != nil {
			log.Print(err)
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		err = ts.ExecuteTemplate(w, "fantasy", data)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}

	default:
		fmt.Println(addId)
		data, err := AllFantasyPlayers()
		if err != nil {
			log.Print(err)
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		err = ts.ExecuteTemplate(w, "fantasy", data)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}

	}
}

func myTeam(w http.ResponseWriter, r *http.Request) {
	res, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	text := res.Query()
	username := text.Get("user")
	data, err := dbfantasy(username)
	fmt.Println(username)
	if err != nil {
		log.Print(err)
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.ExecuteTemplate(w, "myteam", data)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}

func addPlayer(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "post" {
		r.ParseForm()
		fmt.Println(r.Form["ID"])
	} else {
		r.ParseForm()
		fmt.Println(r.Form["ID"])

		// logic part of log in
	}
}
