package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

/*
TODO: Pull individual meetup resource link
* verify that the document body contains the appropriate div resource
* Find the right class
* curate a list of upcoming meetups
*/
func main() {
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	resp, err := get(client, "https://www.meetup.com/pro/forge-utah/")
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	urls, err := GetMeetups(body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(urls)

}
func GetMeetups(body []byte) ([]string, error) {
	var urls []string

	// split in to "script" elements
	str2 := strings.SplitAfter(string(body), `</script>`)
	for _, s := range str2 {
		// search for url key
		if strings.Contains(s, `"url"`) {
			// clean up html
			m := make(map[string]interface{})
			meetupInfo := cleanHTMLReactScriptTag(s)
			err := json.Unmarshal([]byte(meetupInfo), &m)
			if err != nil {
				return []string{}, err
			}
			urls = append(urls, m["url"].(string))
		}
	}
	return urls, nil
}

func cleanHTMLReactScriptTag(s string) string {
	s = strings.TrimPrefix(s, `<script data-react-helmet="true" type="application/ld+json">`)
	s = strings.TrimSuffix(s, `</script>`)
	return s
}

func get(client *http.Client, url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.AddCookie(&http.Cookie{
		Name:  "name",
		Value: "value",
	})
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("not a 200 status code")
	}
	return resp, err
}
