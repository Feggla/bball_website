package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
    host = "localhost"
    user = "newuser"
    password = "password"
    dbname = "postgres"
)

func main() {
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

	// sqlStatement := `INSERT INTO users (id, username)
	// Values ($1, $2)`
	// _, err = db.Exec(sqlStatement, 200, "michaelfeggans@mail.com")
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Println("Added successfully")
	// }
	var (
		userID       int64
		name string
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
			id:userID,
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