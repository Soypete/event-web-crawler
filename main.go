package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Soypete/event-web-crawler/firebase"
	"github.com/Soypete/event-web-crawler/meetup"
)

/*
TODO:
- deploy script to run weekly
- update site with githib api
*/
func handler(w http.ResponseWriter, r *http.Request) {
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
	// take infos and store firebase
	err = firebaseClt.AddMeetupInfos(ctx, "Meetups", infos)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Saved meetups in firestore\n")
}

func main() {
	log.Print("starting server...")
	http.HandleFunc("/crawl", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
