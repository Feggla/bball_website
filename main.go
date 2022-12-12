package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	// "net/url"
	"os"
)

func main() {
	var page int
	page = 1
	// params := url.Values{}
	// params.Add("page", page)
	for i := 0; i <155; i++{
		url := fmt.Sprint("https://www.balldontlie.io/api/v1/players?page=", page)
		// resp, err := http.Get("https://www.balldontlie.io/api/v1/players"+params.Encode())
		resp, err := http.Get(url)
		if err != nil {
		log.Fatalln(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
		log.Fatalln(err)
		}

		// sb := string(body)
		// fmt.Printf("%T", sb)
		// fmt.Print(sb)


		type Player struct {
			Data [] struct {
				ID int `json:"id`
				First_name string `json:"first_name"`
				Height interface{} `json:"height_feet"`
				Height_inches interface{} `json:"heigh_inches"`
				Last_name string `json:"last_name"`
				Position string `json:"position"`
				Team struct {
					ID int `json:"id"`
					Abbreviation string `json:"abbreviation"`
					City string `json:"city"`
					Conference string `json:"conference"`
					Division string `json:"division"`
					Full_name string `json:"full_name"`
					Name string `json:"name"`
				} `json:"team`
				Weight interface{} `json:"weight_pounds"`
			} `json:"data"`
			Meta [] struct {
				Pages int `json:"total_pages"`
				Current_page int `json:"current_page"`
				Next_page int `json:"next_page"`
				Per_page int `json:"per_page"`
				Total_count int `json:"total_count"`
			} `json:"meta"`
			}
		
		
		var m Player
		err = json.Unmarshal(body, &m)
		// for i := 0; i < len(m.Data); i++ {
		// 	fmt.Println(m.Data[i].Position)

		// }
		if len(os.Args) == 1 {
			fmt.Println("Please select a basketball pos: G for Guard, C for Centre of F for Forward")
			break
		} else {
			pos := os.Args[1:][0]
		for i := 0; i<len(m.Data); i++ {
			if m.Data[i].Position == strings.ToUpper(pos) {
				fmt.Println(m.Data[i].First_name, m.Data[i].Last_name, m.Data[i].Position)
				fmt.Println(m.Data[i].Team.Abbreviation)
			}
		}
		// pos := os.Args[1:][0]
		// for i := 0; i<len(m.Data); i++ {
		// 	if m.Data[i].Position == strings.ToUpper(pos) {
		// 		fmt.Println(m.Data[i].First_name, m.Data[i].Last_name, m.Data[i].Position)
		// 		fmt.Println(m.Data[i].Team.Abbreviation)
			}
		// }
		page += 1

	}
}	