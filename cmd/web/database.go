package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("host")
	user     = os.Getenv("user")
	password = os.Getenv("password")
	dbname   = os.Getenv("dbname")
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

func dbfantasy(username string) ([]Fantasydb, error) {
	var fantasy []Fantasydb

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbname, password, host)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return []Fantasydb{}, err
	}
	defer db.Close()

	a := `
	SELECT users.username, player.first_name, player.last_name, player.position, player.team, player.player_id
	FROM fantasy
	JOIN users ON users.id = fantasy.user_id
	JOIN player ON player.player_id = fantasy.player_id WHERE users.username = $1`

	// querystring := fmt.Sprintf("SELECT users.username, player.first_name, player.last_name, player.position, player.team, player.player_id  FROM fantasy JOIN users ON users.id = fantasy.user_id JOIN player ON player.player_id = fantasy.player_id WHERE users.username = '%s'", query)
	rows, err := db.Query(a, username)
	if err != nil {
		return []Fantasydb{}, err
	}
	defer rows.Close()
	var fantdata Fantasydb
	for rows.Next() {
		err := rows.Scan(&fantdata.Player.FantasyPlayer, &fantdata.Player.FirstName, &fantdata.Player.LastName, &fantdata.Player.Position, &fantdata.Player.Team, &fantdata.Player.Id)
		if err != nil {
			return []Fantasydb{}, err
		}
		fantasy = append(fantasy, fantdata)
	}
	return fantasy, nil
}

func AllFantasyPlayers(username string) ([]Fantasydb, error) {
	var fantasy []Fantasydb
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbname, password, host)
	db, err := sql.Open("postgres", connStr)
	fmt.Println("Working")
	if err != nil {
		return []Fantasydb{}, err
	}
	defer db.Close()
	querystring := fmt.Sprintf("SELECT player.first_name, player.last_name, player.position, player.team, player.player_id FROM player WHERE player.player_id NOT IN (SELECT fantasy.player_id FROM fantasy JOIN users on users.id = fantasy.user_id WHERE users.username = '%s')", username)
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
	sqlstatement := "INSERT INTO player (first_name, last_name, position, team, player_id) VALUES ($1, $2, $3, $4, $5)"
	for _, p := range all_players {
		_, err := db.Exec(sqlstatement, p.First_name, p.Last_name, p.Position, p.Team.Abbreviation, p.ID)
		if err != nil {
			panic(err)
		}
		name := fmt.Sprintf("Player added: %s %s", p.First_name, p.Last_name)
		fmt.Print(name)
	}
}

func dbCheckLog(userid string) (string, error) {
	var usernames []string
	var name string
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbname, password, host)
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	querystring := "SELECT username FROM users"
	rows, err := db.Query(querystring)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name)
		usernames = append(usernames, name)
		if err != nil {
			return "", err
		}
	}
	for _, x := range usernames {
		if x == userid {
			return x, nil
		}
	}
	return "", nil
}

func addPlayer(playerid int, username string) error {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbname, password, host)
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	// querystring := fmt.Sprintf("INSERT INTO fantasy (user, player) SELECT users.username, player.id FROM (VALUES ('%s', %d)) JOIN users USING (username) JOIN player USING (id)", username, playerid)
	q := fmt.Sprintf("INSERT INTO fantasy SELECT u.id, p.player_id FROM users u, player p WHERE u.username = '%s' and p.player_id = %d ;", username, playerid)

	_, err = db.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

func removePlayer(playerid int, username string) error {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", user, dbname, password, host)
	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	q := fmt.Sprintf("DELETE FROM fantasy USING users WHERE users.id = fantasy.user_id AND users.username = '%s' AND fantasy.player_id = %d ;", username, playerid)
	x, err := db.Exec(q)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(x)
	return nil
}
