package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type MeetupInfo struct {
	Name                string `json:"name"`
	URL                 string `json:"url"`
	Description         string `json:"description"`
	Startdate           string `json:"startDate"`
	Enddate             string `json:"endDate"`
	Eventstatus         string `json:"eventStatus"`
	Eventattendancemode string `json:"eventAttendanceMode"`
	Location            struct {
		Type string `json:"@type"`
		URL  string `json:"url"`
	} `json:"location"`
}

/*
TODO:
- pull meetup information
	- Title, description, date, location (online/in-person)
- store information - as json file??
- add github actions for testing
- clean up code and pull out into directory
- deploy script to run weekly
*/

func main() {
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	file, err := os.OpenFile("datums/meetups.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)

	resp, err := get(client, "https://www.meetup.com/pro/forge-utah/")
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	urls, err := GetMeetupsURLs(body)
	if err != nil {
		log.Fatal(err)
	}
	var infos []MeetupInfo
	for i, url := range urls {
		info, err := GetMeetupInfo(client, i, url)
		if err != nil {
			log.Fatal(err)
		}
		infos = append(infos, info)
	}
	fmt.Println(len(infos))
	encoder.Encode(&infos)
	fmt.Println("done")
}

func GetMeetupInfo(client *http.Client, count int, url string) (MeetupInfo, error) {
	resp, err := get(client, url)
	if err != nil {
		return MeetupInfo{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return MeetupInfo{}, err
	}
	defer resp.Body.Close()
	var parsedMeetup MeetupInfo
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
			fmt.Println(parsedMeetup.URL)
			continue
		}
	}
	return parsedMeetup, nil
}

// GetMeetupURLs pulls a list of upcoming Forge Foundation
func GetMeetupsURLs(body []byte) ([]string, error) {
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
