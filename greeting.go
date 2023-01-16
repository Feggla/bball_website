package main

import (
	"fmt"
	"os"
	"os"
	"time"
)

func main() {
	name := os.Args[1:][0]
	var time_of_day string
	dt := time.Now()
	current_time := dt.Format("15:04:05")
	now := time.Now().Hour()
	if now >= 17 && now < 20 {
		time_of_day = "evening"
	} else if now < 17 && now > 7 {
		time_of_day = "daytime"
	} else if now > 20 && now <= 7 {
		time_of_day = "night time"
	}
	fmt.Print("Hey, ", name, ". The time is ", current_time, " and it is ", time_of_day, ".")
}
