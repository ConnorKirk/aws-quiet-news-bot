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

func containsRegion(item *gofeed.Item, checkRegion string) bool {
	return strings.Contains(strings.ToLower(item.Title), strings.ToLower(regions[checkRegion]))
}

// excludeRegion takes a feed containing a list of items
// If an item contains an excluded region, it is removed from the list
func excludeRegion(region string) filterFunc {
	return func(f *gofeed.Feed) {
		var newItems []*gofeed.Item
		for _, item := range f.Items {
			if !containsRegion(item, region) {
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

			// Dont exclude kept regions
			for _, keepRegion := range keepRegions {
				if regions[keepRegion] == r {
					continue excludeLoop
				}
			}

			excludeRegion(r)(f)
		}
	}
}

func getRegions(item *gofeed.Item) (regions []string, ok bool) {
	// ok itialises to false
	ok = false
	for _, region := range regions {
		if strings.Contains(strings.ToLower(item.Description), strings.ToLower(region)) {
			regions = append(regions, region)
			ok = true
		}
	}

	return regions, ok
}
