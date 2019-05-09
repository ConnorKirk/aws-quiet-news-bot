package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/lunny/html2md"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mmcdole/gofeed"
)

const (
	url        = "https://aws.amazon.com/new/feed/"
	ENVwebhook = "webhook"
)

var (
	timeWindow = 24 * time.Hour

	// Destination Webhook
	webhookURL string
)

func init() {
	webhookURL = os.Getenv(ENVwebhook)
	log.Printf("Using %v", webhookURL)
}

func handler() error {
	// get RSS feed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		log.Printf("Error ParseURL(%v): %v", url, err)
		return err
	}

	// filter RSS feed
	opts := []filterFunc{filterDate(timeWindow)}
	for _, opt := range opts {
		opt(feed)
	}

	if len(feed.Items) == 0 {
		log.Printf("No items to display")
	}

	sort.Slice(feed.Items,
		func(i, j int) bool {
			return feed.Items[i].PublishedParsed.Sub(*feed.Items[j].PublishedParsed) < 0
		})

	// Process Feed Items
	client := http.Client{}
	for _, item := range feed.Items {
		err = postItem(client, item)

		// Respect the rate limit
		time.Sleep(1 * time.Second)

	}
	return err
}

func main() {
	lambda.Start(handler)
}

func postItem(client http.Client, item *gofeed.Item) error {
	c := struct {
		Content string
	}{
		Content: buildContent(item),
	}

	b, err := json.Marshal(c)
	if err != nil {
		log.Printf("ERROR json.Marshal(): %v", err)
		return err
	}

	// Buil Request
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		log.Printf("Error http.NewRequest(%b): %v", b, err)
		return err
	}

	// Send Request
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error client.Do(): %v", err)
		return err
	}
	if resp.StatusCode != 200 {
		log.Printf("err client.Do(): non 200 status code: %v - %s", resp.StatusCode, resp.Status)
		return err
	}

	return nil
}

type filterFunc func(*gofeed.Feed)

// filterDate excludes all items older than the provided duration
func filterDate(d time.Duration) filterFunc {
	return func(f *gofeed.Feed) {
		var newItems []*gofeed.Item
		for _, i := range f.Items {
			published := *i.PublishedParsed
			now := time.Now()
			if now.Sub(published) > d {
				continue
			}
			newItems = append(newItems, i)
		}
		f.Items = newItems
	}
}

func containsRegion(item *gofeed.Item, checkRegion []string) bool {
	for _, r := range checkRegion {
		if strings.Contains(item.Title, regions[r]) {
			return true
		}
	}
	return false
}

// excludeRegion takes a feed containing a list of items
// If an item contains an excluded region, it is removed from the list
func excludeRegion(regions ...string) filterFunc {
	return func(f *gofeed.Feed) {
		var newItems []*gofeed.Item
		for _, item := range f.Items {
			if !containsRegion(item, regions) {
				newItems = append(newItems, item)
			}
		}
		f.Items = newItems
	}
}

// excludeAllExcept removes
func excludeAllExcept(keepRegions ...string) filterFunc {
	return func(f *gofeed.Feed) {

	excludeLoop:

		for _, r := range regions {

			// Skip excluding kept regions
			for _, keepRegion := range keepRegions {
				if regions[keepRegion] == r {
					continue excludeLoop
				}
			}

			excludeRegion(r)(f)
		}
	}
}

// buildContent takes an item and returns a markdown formatted
// string containing the publish date, title and description
func buildContent(item *gofeed.Item) string {
	var b bytes.Buffer

	// Signal markdown content
	b.WriteString("/md ")

	const format = "2006-01-02"
	b.WriteString(italics(item.PublishedParsed.Format(format)) + "\n")
	b.WriteString(link(item.Title, item.Link) + "\n")

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
