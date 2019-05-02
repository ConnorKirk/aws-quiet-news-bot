package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lunny/html2md"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mmcdole/gofeed"
)

const (
	url     = "https://aws.amazon.com/new/feed/"
	webhook = "https://hooks.chime.aws/incomingwebhooks/52bcd653-234a-4f19-ace7-06dab9238d15?token=RXB0aUtlemZ8MXxLSGNRUERqLUs1MVBYbzQwLUNnQ3hYR2VJcXNUbThFemt3SUFuM25iazlr"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

func handler() error {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Printf("Error ParseURL(%v): %v", url, err)
		return err
	}

	opts := []filterFunc{filterDate}
	for _, opt := range opts {
		feed = opt(feed)
	}

	client := http.Client{}
	for _, item := range feed.Items {
		c := post{
			Content: buildContent(item),
		}

		b, err := json.Marshal(c)
		if err != nil {
			log.Printf("ERROR json.Marshal(): %v", err)
			return err
		}
		fmt.Println(string(b))
		req, err := http.NewRequest(http.MethodPost, webhook, bytes.NewReader(b))
		req.Header.Add("Content-Type", "application/json")
		if err != nil {
			log.Printf("Error http.NewRequest(%b): %v", b, err)
			return err
		}
		client.Do(req)

	}
	return nil
}

func main() {
	lambda.Start(handler)
}

type filterFunc func(*gofeed.Feed) *gofeed.Feed

func filterDate(f *gofeed.Feed) *gofeed.Feed {
	var newItems []*gofeed.Item
	for _, i := range f.Items {
		published := *i.PublishedParsed
		now := time.Now()
		if now.Sub(published) > 24*time.Hour {
			break
		}
		newItems = append(newItems, i)
	}
	f.Items = newItems
	return f
}
func containsRegion(item *gofeed.Item, regions []string) bool {
	for _, r := range regions {
		if strings.Contains(item.Title, r) {
			return true
		}
	}
	return false
}

func excludeRegion(regions ...string) filterFunc {
	return func(f *gofeed.Feed) *gofeed.Feed {
		var newItems []*gofeed.Item
		for _, item := range f.Items {
			if !containsRegion(item, regions) {
				newItems = append(newItems, item)
			}
		}
		f.Items = newItems
		return f
	}
}

func includeRegion(region ...string) filterFunc {
	return func(f *gofeed.Feed) *gofeed.Feed {
		var newItems []*gofeed.Item

		for _, item := range f.Items {
			if containsRegion(item, region) {
				newItems = append(newItems, item)
			}
		}
		f.Items = newItems
		return f
	}
}

type post struct {
	Content string
}

// buildContent takes an item and returns a markdown formatted
// string containing the publish date, title and description
func buildContent(item *gofeed.Item) string {
	var b bytes.Buffer

	// Signal markdown content
	b.WriteString("/md ")

	b.WriteString(italics(item.Published) + "\n")
	b.WriteString(link(item.Title, item.Link) + "\n\n")

	desc := html2md.Convert(item.Description)
	b.WriteString(desc)
	return b.String()
}

// Convert to markdown italics
func italics(s string) string {
	return fmt.Sprintf("*%s*", s)
}

// Build a markdown formatted link
func link(s string, link string) string {
	return fmt.Sprintf("[%s](%s)", s, link)
}
