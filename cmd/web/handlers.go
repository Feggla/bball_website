package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"os"
	"strconv"

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
	"./ui/html/pages/myteamtable.tmpl.html",
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
	var (
		host     = os.Getenv("host")
		user     = os.Getenv("user")
		password = os.Getenv("password")
		dbname   = os.Getenv("dbname")
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
		passw   string
	)
	type userind struct {
		Id       int64
		Username string
		Password string
	}
	rows, err := db.Query("SELECT id, username, password FROM users order by id asc")
	if err != nil {
		panic(err)
	}
	var users []userind
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&userID, &name, &passw)
		if err != nil {
			panic(err)
		}
		users = append(users, userind{
			Id:       userID,
			Username: name,
			Password: passw,
		})
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
	cookie, err := r.Cookie("user")
	if err != nil {
		username := ""
		data, err := AllFantasyPlayers(username)
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
	username := cookie.Value
	err = r.ParseForm()
	if err != nil {
		log.Print(err)
	}
	switch r.Method {
	case "POST":
		id := r.PostForm["addID"]
		fmt.Println(id)
		addId, err := strconv.Atoi(id[0])
		if err != nil {
			log.Print(err)
		}
		err = addPlayer(addId, username)
		fmt.Println(err)
		data, err := AllFantasyPlayers(username)
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
	case "GET":
		data, err := AllFantasyPlayers(username)
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
	err := r.ParseForm()
	if err != nil {
		log.Print(err)
	}
	cookie, err := r.Cookie("user")
	if err != nil {
		res, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		text := res.Query()
		username := text.Get("user")
		pass := text.Get("pass_input")
		name, err := dbCheckLog(username, pass)
		if err != nil {
			log.Print(err)
		}
		fmt.Println(username)

		switch name {
		case "":
			fmt.Println("Username not found in db")
			ts, err := template.ParseFiles(files...)
			if err != nil {
				log.Print(err.Error())
				http.Error(w, "Internal Server Error", 500)
				return
			}
			err = ts.ExecuteTemplate(w, "myteam", nil)
			if err != nil {
				log.Print(err.Error())
				http.Error(w, "Internal Server Error", 500)
			}
		default:
			switch r.Method {
			case "POST":
				id := r.PostForm["removeid"]
				fmt.Println(id)
				addId, err := strconv.Atoi(id[0])
				if err != nil {
					log.Print(err)
				}
				err = removePlayer(addId, username)
				if err != nil {
					log.Print(err)
				}
				cookie := http.Cookie{Name: "user", Value: name}
				http.SetCookie(w, &cookie)
				data, err := dbfantasy(cookie.Value)
				if err != nil {
					log.Print(err)
				}
				fmt.Println(err)
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
			case "GET":
				cookie := http.Cookie{Name: "user", Value: name}
				http.SetCookie(w, &cookie)
				data, err := dbfantasy(cookie.Value)
				if err != nil {
					log.Print(err)
				}
				fmt.Println(err)
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
		}
	} else {
		res, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		text := res.Query()
		username := text.Get("user")
		password := text.Get("pass_input")
		name, _ := dbCheckLog(username, password)
		if name == "" {
			username = cookie.Value
		}
		switch r.Method {
		case "POST":
			id := r.PostForm["removeid"]
			fmt.Println(id)
			addId, err := strconv.Atoi(id[0])
			if err != nil {
				log.Print(err)
			}
			err = removePlayer(addId, username)
			if err != nil {
				log.Print(err)
			}
			cookie := http.Cookie{Name: "user", Value: username}
			http.SetCookie(w, &cookie)
			data, err := dbfantasy(cookie.Value)
			if err != nil {
				log.Print(err)
			}
			fmt.Println(err)
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
		case "GET":
			cookie := http.Cookie{Name: "user", Value: username}
			http.SetCookie(w, &cookie)
			data, err := dbfantasy(cookie.Value)
			if err != nil {
				log.Print(err)
			}
			fmt.Println(err)
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
	}
}
