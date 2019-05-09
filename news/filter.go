package main

import (
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

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
