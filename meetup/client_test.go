package meetup

import (
	"net/http"
	"reflect"
	"testing"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		name string
		want *Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Setup(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Setup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetProPage(t *testing.T) {
	tests := []struct {
		name    string
		c       *Client
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.GetProPage()
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetProPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetProPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_get(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		want    *http.Response
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.get(tt.args.url)
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

func TestClient_GetMeetupInfo(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		want    Info
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.GetMeetupInfo(tt.args.url)
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
