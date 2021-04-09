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

// Client is meetup wrapper for http.Client
type Client struct {
	client *http.Client
	proURL string
}

// Setup the meetup http client
func Setup() *Client {
	c := Client{
		client: &http.Client{
			Timeout: time.Second * 5,
		},
		proURL: "https://www.meetup.com/pro/forge-utah/",
	}
	return &c
}

// GetProPage retrieves the HTML payload from meetups ForgeFoundation pro page
func (c *Client) GetProPage() ([]byte, error) {
	// check http.Client initialized
	if c.client == nil {
		return []byte{}, errors.New("http.Client not initialized")
	}
	resp, err := c.GetWebPage(c.proURL)
	if err != nil {
		fmt.Println(resp)
		return []byte{}, err
	}
	return c.GetWebPage(c.proURL)
}

func (c *Client) GetWebPage(url string) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Cannot create request %w", err)
	}
	request.AddCookie(&http.Cookie{
		Name:  "name",
		Value: "value",
	})
	request.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(request)
	if err != nil {
		return []byte{}, fmt.Errorf("failure in Do request:\n %w ---\n", err)
	}
	if resp.StatusCode != http.StatusOK {
		return []byte{}, errors.New("not a 200 status code")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Cannot parse response body: %w", err)
	}
	defer resp.Body.Close()
	return body, nil
}

func (c *Client) GetMeetupInfo(url string) (Info, error) {
	body, err := c.GetWebPage(url)
	if err != nil {
		return Info{}, err
	}
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
