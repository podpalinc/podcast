package parser_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	parser "github.com/podpalinc/rss-feed-generator/parser"
	"github.com/stretchr/testify/assert"
)

func TestDetectFeedType(t *testing.T) {
	var feedTypeTests = []struct {
		file     string
		expected parser.FeedType
	}{
		{"atom03_feed.xml", parser.FeedTypeAtom},
		{"atom10_feed.xml", parser.FeedTypeAtom},
		{"rss_feed.xml", parser.FeedTypeRSS},
		{"rss_feed_bom.xml", parser.FeedTypeRSS},
		{"rss_feed_leading_spaces.xml", parser.FeedTypeRSS},
		{"rdf_feed.xml", parser.FeedTypeRSS},
		{"unknown_feed.xml", parser.FeedTypeUnknown},
		{"empty_feed.xml", parser.FeedTypeUnknown},
		{"json10_feed.json", parser.FeedTypeJSON},
	}

	for _, test := range feedTypeTests {
		fmt.Printf("Testing %s... ", test.file)

		// Get feed content
		path := fmt.Sprintf("testdata/parser/universal/%s", test.file)
		f, _ := ioutil.ReadFile(path)

		// Get actual value
		actual := parser.DetectFeedType(bytes.NewReader(f))

		if assert.Equal(t, actual, test.expected, "Feed file %s did not match expected type %d", test.file, test.expected) {
			fmt.Printf("OK\n")
		} else {
			fmt.Printf("Failed\n")
		}
	}
}

// Examples

func ExampleDetectFeedType() {
	feedData := `<rss version="2.0">
<channel>
<title>Sample Feed</title>
</channel>
</rss>`
	feedType := parser.DetectFeedType(strings.NewReader(feedData))
	if feedType == parser.FeedTypeRSS {
		fmt.Println("Wow! This is an RSS feed!")
	}
}
