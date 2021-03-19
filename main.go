package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	in := make(chan string, 10)
	for i := 0; i < 5; i++ {
		go request(client, in)
	}

	resp, err := get(client, "https://www.meetup.com/pro/forge-utah/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}
	selection := doc.Find("a")
	if len(selection.Nodes) == 0 {
		panic("no 'a' nodes found")
	}
	for _, node := range selection.Nodes {
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				if strings.HasPrefix(attr.Val, "http") {
					in <- attr.Val
				}
				break
			}
		}
	}
}

func request(client *http.Client, in chan string) {
	for {
		select {
		case url := <-in:
			err := getChild(client, url)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func getChild(client *http.Client, attr string) error {
	childResp, err := get(client, attr)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(childResp.Body)
	if err != nil {
		return err
	}
	defer childResp.Body.Close()
	if childResp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s is not status 200, status: %d", attr, childResp.StatusCode)
	}
	length := len(b)
	fmt.Println("URL :", attr, "Status Code :", childResp.StatusCode, "length:", length)
	return nil
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
		return nil, fmt.Errorf("not a 200 status code")
	}
	return resp, err
}
