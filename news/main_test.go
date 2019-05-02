package main

import (
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestIncludeRegion(t *testing.T) {
	feed := newFeed()
	blameCanada := includeRegion(cacentral1)
	feed = blameCanada(feed)
	got := len(feed.Items)
	want := 1
	if got != want {
		t.Errorf("Expected Length()=%v; got=%v", want, got)
	}
}

func TestExcludeRegion(t *testing.T) {
	feed := newFeed()
	blameCanada := excludeRegion(cacentral1)
	feed = blameCanada(feed)
	got := len(feed.Items)
	want := 0
	if got != want {
		t.Errorf("Expected Length()=%v; got=%v", want, got)
	}
}

func newFeed() *gofeed.Feed {
	i := &gofeed.Item{
		Title: "Blame Canada",
	}

	return &gofeed.Feed{
		Items: []*gofeed.Item{
			i,
		},
	}
}
