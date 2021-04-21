package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/Soypete/event-web-crawler/meetup"
)

/*
TODO:
- deploy script to run weekly
*/

func main() {
	meetupClt := meetup.Setup()
	file, err := os.OpenFile("datums/meetups.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// HTMLfile, err := os.OpenFile("datums/meetup.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)

	body, err := meetupClt.GetProPage()
	if err != nil {
		log.Fatal(err)
	}
	// HTMLfile.WriteString(string(body))
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
	encoder.Encode(&infos)
}
