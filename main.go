package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Soypete/event-web-crawler/firebase"
	"github.com/Soypete/event-web-crawler/meetup"
)

/*
TODO:
- deploy script to run weekly
*/

func main() {
	ctx := context.Background()
	firebaseClt, err := firebase.Setup(ctx)
	if err != nil {
		log.Fatal(err)
	}
	meetupClt := meetup.Setup()
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
	data, err := json.Marshal(infos)
	if err != nil {
		log.Fatal(err)
	}
	// take infos and store firebase
	err = firebaseClt.AddToCollection(ctx, "Meetups", data)
	if err != nil {
		log.Fatal(err)
	}
	// push to site via github api
}
