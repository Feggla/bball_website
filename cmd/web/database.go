package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	user     = "newuser"
	password = "password"
	dbname   = "postgres"
)

type Fantasydb struct {
	Player struct {
		FantasyPlayer string
		FirstName     string
		LastName      string
		Position      string
		Team          string
		Id            int
	}
}

func Dbread() {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbname, password, host)
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
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
		id       int64
		username string
	}
	rows, err := db.Query("SELECT id, username FROM users order by id desc")
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
			id:       userID,
			username: name,
		})

		fmt.Println(users)
		fmt.Println("\n", userID, name)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

}

func dbfantasy(query string) ([]Fantasydb, error) {
	var fantasy []Fantasydb

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbname, password, host)
	db, err := sql.Open("postgres", connStr)
	fmt.Println("Working")
	if err != nil {
		return []Fantasydb{}, err
	}
	defer db.Close()
	querystring := fmt.Sprintf("SELECT users.username, player.first_name, player.last_name, player.position, player.team, player.id  FROM fantasy JOIN users ON users.id = fantasy.user JOIN player ON player.id = fantasy.player WHERE users.username = '%s'", query)
	fmt.Println(querystring)
	rows, err := db.Query(querystring)
	if err != nil {
		return []Fantasydb{}, err
	}
	defer rows.Close()
	var fantdata Fantasydb
	for rows.Next() {
		fmt.Println(rows)
		err := rows.Scan(&fantdata.Player.FantasyPlayer, &fantdata.Player.FirstName, &fantdata.Player.LastName, &fantdata.Player.Position, &fantdata.Player.Team, &fantdata.Player.Id)
		if err != nil {
			return []Fantasydb{}, err
		}
		fantasy = append(fantasy, fantdata)
	}
	return fantasy, nil
}

func AllFantasyPlayers() ([]Fantasydb, error) {
	var fantasy []Fantasydb
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbname, password, host)
	db, err := sql.Open("postgres", connStr)
	fmt.Println("Working")
	if err != nil {
		return []Fantasydb{}, err
	}
	defer db.Close()
	querystring := "SELECT player.first_name, player.last_name, player.position, player.team, player.id  FROM player"
	rows, err := db.Query(querystring)
	if err != nil {
		return []Fantasydb{}, err

	}
	defer rows.Close()
	var fantdata Fantasydb
	for rows.Next() {
		err := rows.Scan(&fantdata.Player.FirstName, &fantdata.Player.LastName, &fantdata.Player.Position, &fantdata.Player.Team, &fantdata.Player.Id)
		if err != nil {
			return []Fantasydb{}, err
		}
		fantasy = append(fantasy, fantdata)
	}
	return fantasy, nil
}
func Dbadd() {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbname, password, host)
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	all_players, err := GetAllPlayers()
	if err != nil {
		log.Println(err)
	}
	sqlstatement := "INSERT INTO player (first_name, last_name, position, team, id) VALUES ($1, $2, $3, $4, $5)"
	for _, p := range all_players {
		_, err := db.Exec(sqlstatement, p.First_name, p.Last_name, p.Position, p.Team.Abbreviation, p.ID)
		if err != nil {
			panic(err)
		}
		name := fmt.Sprintf("Player added: %s %s", p.First_name, p.Last_name)
		fmt.Print(name)
	}
}
