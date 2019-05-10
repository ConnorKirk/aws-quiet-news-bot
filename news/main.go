package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/lunny/html2md"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mmcdole/gofeed"
)

const (
	// Viper Keys
	vPrefix           = "news"
	vInputFeedURL     = "INPUT_FEED_URL"
	vOutputWebHookURL = "OUTPUT_WEBHOOK_URL"
	vTimeWindowDays   = "TIME_WINDOW_DAYS"

	// Defaults
	defaultInputFeedURL   = "https://aws.amazon.com/new/feed/"
	defaultTimeWindowDays = 1
)

var (
	webhookURL     string        // URL to post to
	feedURL        string        // URL to get RSS feed from
	timeWindowDays int           // Number of days in the window
	window         time.Duration // Length of window in seconds
)

func init() {
	viper.SetEnvPrefix(vPrefix)
	viper.SetDefault(vInputFeedURL, defaultInputFeedURL)
	viper.SetDefault(vTimeWindowDays, defaultTimeWindowDays)

	err := viper.BindEnv(vInputFeedURL)
	if err != nil {
		panic(err)
	}
	err = viper.BindEnv(vOutputWebHookURL)
	if err != nil {
		panic(err)
	}
	err = viper.BindEnv(vTimeWindowDays)
	if err != nil {
		panic(err)
	}
	webhookURL = viper.GetString(vOutputWebHookURL)
	feedURL = viper.GetString(vInputFeedURL)
	timeWindowDays = viper.GetInt(vTimeWindowDays)
	window = time.Duration(timeWindowDays) * 24 * time.Hour
}

func handler() error {
	// get RSS feed
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		log.Printf("Error ParseURL(%v): %v", feedURL, err)
		return err
	}

	// filter RSS feed
	opts := []filterFunc{filterDate(window)}
	for _, opt := range opts {
		opt(feed)
	}

	if len(feed.Items) == 0 {
		log.Printf("No items to display")
	}

	// Sort by Published Date
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

// postItem formats an RSS item and posts it to the chime webhook
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
	b.WriteString(strings.TrimSpace(desc))
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
