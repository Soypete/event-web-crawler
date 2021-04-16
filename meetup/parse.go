package meetup

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Info contains relevant meetup information that is to be stored.
type Info struct {
	Name                string   `json:"name"`
	URL                 string   `json:"url"`
	Description         string   `json:"description"`
	Startdate           string   `json:"startDate"`
	Enddate             string   `json:"endDate"`
	Eventstatus         string   `json:"eventStatus"`
	Eventattendancemode string   `json:"eventAttendanceMode"`
	Location            Location `json:"location"`
}

type Location struct {
	Type string `json:"@type"`
	URL  string `json:"url"`
}

func cleanHTMLReactScriptTag(s string) string {
	s = strings.TrimPrefix(s, `<script data-react-helmet="true" type="application/ld+json">`)
	s = strings.TrimSuffix(s, `</script>`)
	return s
}

// GetMeetupURLs pulls a list of upcoming Forge Foundation.
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
				return []string{}, fmt.Errorf("Cannot unmarshal meetup urls", err)
			}
			urls = append(urls, m["url"].(string))
		}
	}
	return urls, nil
}
