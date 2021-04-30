package firebase

import (
	"context"
	"fmt"
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
func Setup(ctx context.Context) (*Client, error) {
	projectID := os.Getenv("PROJECT_ID")
	// Use the application default credentials
	conf := firego.Config{ProjectID: projectID}
	app, err := firego.NewApp(ctx, &conf)
	if err != nil {
		return nil, fmt.Errorf("cannot configure firestore app: %w", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot configure firestore client: %w", err)
	}
	return &Client{
		App:    app,
		Client: client,
	}, nil
}

// AddToCollection adds data to a collection on cloud firestore.
func (c *Client) AddToCollection(ctx context.Context, collectionName string, data interface{}) error {
	dataRef, WriteResult, err := c.Client.Collection(collectionName).Add(ctx, data)
	if err != nil {
		return fmt.Errorf("cannot write to collection %s: %w", collectionName, err)
	}
	log.Println(dataRef, WriteResult)

	return nil
}
