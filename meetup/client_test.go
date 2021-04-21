package meetup

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

var htmlBlob []byte

func setupTestServer() (Client, *httptest.Server) {
	htmlFile, _ := os.Open("../datums/meetup.txt")
	htmlBlob, _ = ioutil.ReadAll(htmlFile)
	defer htmlFile.Close()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(htmlBlob)
	}))
	testClient := Client{
		client: &http.Client{Timeout: time.Second * 5},
		proURL: ts.URL,
	}

	return testClient, ts
}
func TestSetup(t *testing.T) {
	testClient, ts := setupTestServer()
	defer ts.Close()
	tests := []struct {
		name string
		want *Client
	}{
		{
			name: "sucess",
			want: &testClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Setup(); got == nil {
				t.Errorf("Setup() = %v, \nwant %v", got, tt.want)
			}
		})
	}
}
func TestClient_GetWebPage(t *testing.T) {
	testClient, ts := setupTestServer()
	defer ts.Close()
	tests := []struct {
		name    string
		c       *Client
		url     string
		want    []byte
		wantErr bool
	}{
		{
			name:    "success run",
			c:       &testClient,
			url:     testClient.proURL,
			want:    htmlBlob,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.GetWebPage(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetProPage(t *testing.T) {
	testClient, ts := setupTestServer()
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{}"))
	}))
	defer ts.Close()
	defer ts2.Close()
	client2 := Client{
		client: &http.Client{Timeout: time.Second * 5},
		proURL: ts2.URL,
	}
	tests := []struct {
		name    string
		c       *Client
		want    []byte
		wantErr bool
	}{
		{
			name:    "success",
			c:       &testClient,
			want:    htmlBlob,
			wantErr: false,
		},
		{
			name:    "no client",
			c:       &Client{},
			want:    []byte{},
			wantErr: true,
		},
		{
			name:    "bad json",
			c:       &client2,
			want:    []byte("{}"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.GetProPage()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetProPage() error = %v, \nwantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetProPage() = %v, \nwant %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetMeetupInfo(t *testing.T) {
	testClient, _ := setupTestServer()
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{}"))
	}))
	defer ts2.Close()
	client2 := Client{
		client: &http.Client{Timeout: time.Second * 5},
		proURL: ts2.URL,
	}
	tests := []struct {
		name    string
		c       *Client
		url     string
		want    Info
		wantErr bool
	}{
		{
			name: "success run",
			c:    &testClient,
			url:  testClient.proURL,
			want: Info{
				Name:                "TBA - looking for speaker",
				URL:                 "https://www.meetup.com/utahgophers/events/rtcxdsycchbgb/",
				Description:         "",
				Startdate:           "2021-05-04T18:30-06:00",
				Enddate:             "2021-05-04T20:00-06:00",
				Eventstatus:         "https://schema.org/EventScheduled",
				Eventattendancemode: "https://schema.org/OfflineEventAttendanceMode",
				Location: Location{
					Type: "Place",
				},
			},
			wantErr: false,
		},
		{
			name:    "no client",
			c:       &Client{},
			url:     "",
			want:    Info{},
			wantErr: true,
		},
		{
			name:    "bad json",
			c:       &client2,
			url:     client2.proURL,
			want:    Info{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.GetMeetupInfo(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetMeetupInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetMeetupInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
