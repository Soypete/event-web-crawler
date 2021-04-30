package firebase

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firego "firebase.google.com/go"
)

// Client stores firestore configured object that are needed
// to save data to firestore.
type Client struct {
	App    *firego.App
	Client *firestore.Client
}

// Setup retrieves the necessary project information to set up
// a firestore client.
func Setup(ctx context.Context) Client {
	projectID := os.Getenv("PROJECT_ID")
	// Use the application default credentials
	conf := &firego.Config{ProjectID: projectID}
	app, err := firego.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return Client{
		App:    app,
		Client: client,
	}
}
