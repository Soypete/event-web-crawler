package main

import (
	"context"
	"fmt"
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
func runCrawler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	firebaseClt, err := firebase.Setup(ctx)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	meetupClt := meetup.Setup()
	body, err := meetupClt.GetProPage()
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	urls, err := meetup.GetMeetupsURLs(body)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	var infos []meetup.Info
	for _, url := range urls {
		info, err := meetupClt.GetMeetupInfo(url)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		infos = append(infos, info)
	}
	fmt.Fprintf(w, "%d meetups found\n", len(infos))
	// take infos and store firebase
	err = firebaseClt.AddMeetupInfos(ctx, "Meetups", infos)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Fprintf(w, "Saved meetups in firestore\n")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("we are live"))
}

func main() {
	log.Print("starting server...")
	http.HandleFunc("/crawl", runCrawler)
	http.HandleFunc("/health", healthCheck)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Fatal("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
