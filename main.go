package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Soypete/event-web-crawler/meetup"
)

/*
TODO:
- clean up code and pull out into directory
- add github actions for testing
- deploy script to run weekly
*/

func main() {
	meetupClt := meetup.Setup()
	file, err := os.OpenFile("datums/meetups.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)

	body, err := meetupClt.GetProPage()
	if err != nil {
		log.Fatal(err)
	}
	urls, err := meetup.GetMeetupsURLs(body)
	if err != nil {
		log.Fatal(err)
	}
	var infos []meetup.Info
	for _, url := range urls {
		info, err := meetupClt.GetMeetupInfo(url)
		if err != nil {
			log.Fatal(err)
		}
		infos = append(infos, info)
	}
	fmt.Println(len(infos))
	encoder.Encode(&infos)
	fmt.Println("done")
}
