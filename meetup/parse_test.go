package meetup

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func Test_cleanHTMLReactScriptTag(t *testing.T) {
	testString := `<script data-react-helmet="true" type="application/ld+json">{"@context":"http://schema.org","@type":"Event","name":"Intro To Convolutional Neural Networks","url":"https://www.meetup.com/Machine-Learning-Utah/events/tmnwdsyccgblc/","description":"https://meet.google.com/mvz-vyxv-kzp\nDescription: Neural networks are great for complex data sets, but some sets have more features to figure out than others. Many times these features are initialized based on heuristics and they have to be tuned as the model returns predictions. With convolutional neural networks, the model tunes the features for itself.\n\nIn this talk, you will learn some use cases for CNNs, how they work under the hood, and how you can create a CNN in Python. You’ll be able to","startDate":"2021-04-28T18:30-06:00","endDate":"2021-04-28T20:03-06:00","eventStatus":"https://schema.org/EventScheduled","eventAttendanceMode":"https://schema.org/OnlineEventAttendanceMode","location":{"@type":"VirtualLocation","url":"https://www.meetup.com/Machine-Learning-Utah/events/tmnwdsyccgblc/"},"offers":{"@type":"Offer","price":"0","priceCurrency":"USD","validFrom":"2020-03-24","availability":"https://schema.org/InStock"},"organizer":{"@type":"Organization","name":"Machine Learning Utah","url":"https://www.meetup.com/Machine-Learning-Utah/"}}</script>`
	wantString := `{"@context":"http://schema.org","@type":"Event","name":"Intro To Convolutional Neural Networks","url":"https://www.meetup.com/Machine-Learning-Utah/events/tmnwdsyccgblc/","description":"https://meet.google.com/mvz-vyxv-kzp\nDescription: Neural networks are great for complex data sets, but some sets have more features to figure out than others. Many times these features are initialized based on heuristics and they have to be tuned as the model returns predictions. With convolutional neural networks, the model tunes the features for itself.\n\nIn this talk, you will learn some use cases for CNNs, how they work under the hood, and how you can create a CNN in Python. You’ll be able to","startDate":"2021-04-28T18:30-06:00","endDate":"2021-04-28T20:03-06:00","eventStatus":"https://schema.org/EventScheduled","eventAttendanceMode":"https://schema.org/OnlineEventAttendanceMode","location":{"@type":"VirtualLocation","url":"https://www.meetup.com/Machine-Learning-Utah/events/tmnwdsyccgblc/"},"offers":{"@type":"Offer","price":"0","priceCurrency":"USD","validFrom":"2020-03-24","availability":"https://schema.org/InStock"},"organizer":{"@type":"Organization","name":"Machine Learning Utah","url":"https://www.meetup.com/Machine-Learning-Utah/"}}`
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "success",
			s:    testString,
			want: wantString,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanHTMLReactScriptTag(tt.s); got != tt.want {
				t.Errorf("cleanHTMLReactScriptTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMeetupsURLs(t *testing.T) {
	htmlFile, _ := os.Open("../datums/meetup.txt")
	htmlBlob, _ = ioutil.ReadAll(htmlFile)
	defer htmlFile.Close()
	tests := []struct {
		name    string
		body    []byte
		want    []string
		wantErr bool
	}{
		{
			name:    "success",
			body:    htmlBlob,
			want:    []string{"https://www.meetup.com/Machine-Learning-Utah/events/tmnwdsyccgblc/", "https://www.meetup.com/Machine-Learning-Utah/events/tmnwdsycchbjc/", "https://www.meetup.com/Women-Who-Go-Utah/events/bkkzdsyccgbkc/", "https://www.meetup.com/Women-Who-Go-Utah/events/277283469/", "https://www.meetup.com/utahgophers/events/rtcxdsycchbgb/"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMeetupsURLs(tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMeetupsURLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMeetupsURLs() = %v, want %v", got, tt.want)
			}
		})
	}
}
