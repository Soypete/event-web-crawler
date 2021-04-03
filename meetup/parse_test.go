package meetup

import (
	"reflect"
	"testing"
)

func Test_cleanHTMLReactScriptTag(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanHTMLReactScriptTag(tt.args.s); got != tt.want {
				t.Errorf("cleanHTMLReactScriptTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMeetupsURLs(t *testing.T) {
	type args struct {
		body []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMeetupsURLs(tt.args.body)
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
