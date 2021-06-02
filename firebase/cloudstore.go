package firebase

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/Soypete/event-web-crawler/meetup"
)

// Client stores firestore configured object that are needed
// to save data to firestore.
type Client struct {
	App    *firebase.App
	Client *firestore.Client
}

const projectID = "meetup-crawler-store"

// Setup retrieves the necessary project information to set up
// a firestore client.
func Setup(ctx context.Context) (*Client, error) {
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("cannot configure firestore app: %w", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot configure firestore client: %w", err)
	}
	c := &Client{
		App:    app,
		Client: client,
	}
	return c, nil
}

// AddMeetupInfos adds data to a collection on cloud firestore.
func (c *Client) AddMeetupInfos(ctx context.Context, collectionName string, data []meetup.Info) error {
	meetupsInfo := c.Client.Collection(collectionName)
	for _, info := range data {
		_, _, err := meetupsInfo.Doc(time.Now().Format("01-02-06")).Collection("events").Add(ctx, map[string]string{
			"name":                info.Name,
			"url":                 info.URL,
			"Description":         info.Description,
			"Startdate":           info.Startdate,
			"Enddate":             info.Enddate,
			"Eventattendancemode": info.Eventattendancemode,
			"Location":            info.Location.Type,
		})
		if err != nil {
			return fmt.Errorf("cannot save to collection %s: %w", collectionName, err)
		}
	}

	return nil
}
