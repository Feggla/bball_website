package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Player struct {
	ID            int         `json:"id`
	First_name    string      `json:"first_name"`
	Height        interface{} `json:"height_feet"`
	Height_inches interface{} `json:"heigh_inches"`
	Last_name     string      `json:"last_name"`
	Position      string      `json:"position"`
	Team          struct {
		ID           int    `json:"id"`
		Abbreviation string `json:"abbreviation"`
		City         string `json:"city"`
		Conference   string `json:"conference"`
		Division     string `json:"division"`
		Full_name    string `json:"full_name"`
		Name         string `json:"name"`
	} `json:"team`
	Weight interface{} `json:"weight_pounds"`
}

func feedPlayers() {
	log_file := "error_log"
	file, err := os.OpenFile(log_file, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	log.SetOutput(file)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.Println("Logging to custom file")

	if len(os.Args) == 1 {
		fmt.Println("Please select a basketball pos: G for Guard, C for Centre of F for Forward")
		return
	} else {
		input := os.Args[1]
		allPlayers, err := GetAllPlayers()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(len(allPlayers))
		for _, p := range allPlayers {
			if p.Position == strings.ToUpper(input) {
				fmt.Printf("%s %s (%s) %s \n", p.First_name, p.Last_name, p.Team.Abbreviation, p.Position)
			}
		}

	}

}

func GetAllPlayers() ([]Player, error) {
	var players []Player
	for i := 1; i < 3; i++ {
		url := fmt.Sprint("https://www.balldontlie.io/api/v1/players?page=", i)
		resp, err := http.Get(url)

		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		type Response struct {
			Data []Player `json:"data"`
			Meta struct {
				Pages        int `json:"total_pages"`
				Current_page int `json:"current_page"`
				Next_page    int `json:"next_page"`
				Per_page     int `json:"per_page"`
				Total_count  int `json:"total_count"`
			} `json:"meta"`
		}

		var res Response
		err = json.Unmarshal(body, &res)
		if err != nil {
			return nil, err
		}
		players = append(players, res.Data...)
	}
	return players, nil
}

func GetPlayersByPosition(position string) ([]Player, error) {
	var rtn []Player
	all_players, err := GetAllPlayers()
	if err != nil {
		log.Println(err)
		return []Player{}, err
	}
	for _, p := range all_players {
		if p.Position == position {
			rtn = append(rtn, p)

		}

	}
	return rtn, nil
}

func PositionFromQuery(query string) string {
	query = strings.ToLower(query)
	if query == "g" || query == "guard" || query == "guards" || query == "guardians" {
		return "G"
	}
	if query == "c" || query == "centre" || query == "centres" || query == "c-f" {
		return "C"
	}
	if query == "f" || query == "forward" || query == "forwards" || query == "p-f" || query == "s-f" {
		return "F"
	}
	return ""
}

func Guards() ([]Player, error) {
	var guards []Player
	all_players, err := GetAllPlayers()
	if err != nil {
		log.Println(err)
	}
	for _, p := range all_players {
		if p.Position == string('G') {
			guards = append(guards, p)

		}

	}
	return guards, nil
}

func Centres() ([]Player, error) {
	var centres []Player
	all_players, err := GetAllPlayers()
	if err != nil {
		log.Println(err)
	}
	for _, p := range all_players {
		if p.Position == string('C') || p.Position == string("C-F") {
			centres = append(centres, p)
		}
	}
	return centres, nil
}
func Forwards() ([]Player, error) {
	var forwards []Player
	all_players, err := GetAllPlayers()
	if err != nil {
		log.Println(err)
	}
	for _, p := range all_players {
		if p.Position == string('F') || p.Position == string("C-F") {
			forwards = append(forwards, p)
		}
	}
	return forwards, nil
}
