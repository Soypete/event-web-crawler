package meetup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const proURL = "https://www.meetup.com/pro/forge-utah/"

// Client is meetup wrapper for http.Client
type Client struct {
	client *http.Client
}

// Setup the meetup http client
func Setup() *Client {
	c := Client{
		client: &http.Client{
			Timeout: time.Second * 5,
		},
	}
	return &c
}

// GetProPage retrieves the HTML payload from meetups ForgeFoundation pro page
func (c *Client) GetProPage() ([]byte, error) {
	resp, err := c.get(proURL)
	if err != nil {
		return []byte{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()
	return body, nil
}

// TODO: this code is not mine. should I change any of this logic?
func (c *Client) get(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.AddCookie(&http.Cookie{
		Name:  "name",
		Value: "value",
	})
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("not a 200 status code")
	}
	return resp, err
}

func (c *Client) GetMeetupInfo(url string) (Info, error) {
	resp, err := c.get(url)
	if err != nil {
		return Info{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Info{}, err
	}
	defer resp.Body.Close()
	// TODO: this is mostly duplicates on GetMeetupURL
	var parsedMeetup Info
	str2 := strings.SplitAfter(string(body), `</script>`)
	for _, s := range str2 {
		// hope this works?
		if strings.Contains(s, `"url"`) {
			infos := cleanHTMLReactScriptTag(s)
			err = json.Unmarshal([]byte(infos), &parsedMeetup)
			if err != nil {
				continue
			}
			if url == parsedMeetup.URL {
				break
			}
			// every once and a while a random url is parse out of the payload
			fmt.Println(parsedMeetup.URL)
			continue
		}
	}
	return parsedMeetup, nil
}
