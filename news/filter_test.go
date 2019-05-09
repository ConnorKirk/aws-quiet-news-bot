package main

import (
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestExcludeRegion(t *testing.T) {
	feed := newFeed()
	blameCanada := excludeRegion("cacentral1")
	blameCanada(feed)
	got := len(feed.Items)
	want := 2
	if got != want {
		t.Errorf("Expected Length()=%v; got=%v", want, got)
	}
}

func TestExcludeRegionExcept(t *testing.T) {
	feed := newFeed()
	excludeAllExcept("euwest1")(feed)
	got := len(feed.Items)
	want := 1
	if got != want {
		t.Errorf("Expected Length()=%v; got=%v", want, got)
	}

}

func newFeed() *gofeed.Feed {
	i := &gofeed.Item{
		Title: "Blame Canada",
	}
	j := &gofeed.Item{
		Title: "New Service in Dublin",
	}
	k := &gofeed.Item{
		Title: "New Service in Ohio",
	}

	return &gofeed.Feed{
		Items: []*gofeed.Item{
			i,
			j,
			k,
		},
	}
}

func TestContainsRegion(t *testing.T) {
	type args struct {
		item        *gofeed.Item
		checkRegion string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test True",
			args: args{
				item:        newFeed().Items[0],
				checkRegion: "cacentral1",
			},
			want: true,
		},
		{
			name: "Test False",
			args: args{
				item:        newFeed().Items[0],
				checkRegion: "euwest2",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsRegion(tt.args.item, tt.args.checkRegion); got != tt.want {
				t.Errorf("containsRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}
