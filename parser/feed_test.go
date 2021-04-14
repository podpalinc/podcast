package parser_test

import (
	"sort"
	"testing"
	"time"

	parser "github.com/podpalinc/rss-feed-generator/parser"
)

func TestFeedSort(t *testing.T) {
	oldestItem := &parser.Item{
		PublishedParsed: &[]time.Time{time.Unix(0, 0)}[0],
	}
	inbetweenItem := &parser.Item{
		PublishedParsed: &[]time.Time{time.Unix(1, 0)}[0],
	}
	newestItem := &parser.Item{
		PublishedParsed: &[]time.Time{time.Unix(2, 0)}[0],
	}

	feed := parser.Feed{
		Items: []*parser.Item{
			newestItem,
			oldestItem,
			inbetweenItem,
		},
	}
	expected := parser.Feed{
		Items: []*parser.Item{
			oldestItem,
			inbetweenItem,
			newestItem,
		},
	}

	sort.Sort(feed)

	for i, item := range feed.Items {
		if !item.PublishedParsed.Equal(
			*expected.Items[i].PublishedParsed,
		) {
			t.Errorf(
				"Item PublishedParsed = %s; want %s",
				item.PublishedParsed,
				expected.Items[i].PublishedParsed,
			)
		}
	}
}
