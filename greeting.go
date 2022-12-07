package main

import (
	"fmt"
	"time"
	"os"
)

func main() {
	name := os.Args[1:][0]
	var time_of_day string
	current_time := time.Kitchen
	now := time.Now().Hour()
	if now >= 17 && now < 20 {
		time_of_day = "evening"
	} else if now < 17 && now > 7 {
		time_of_day = "daytime"
	} else if now > 20 && now <= 7 {
		time_of_day = "night time"
	}
	fmt.Print("This is an additoin to the feature branch")
	fmt.Print("Hey, ",name,". The time is ", current_time, " and it is ", time_of_day, ".")
}