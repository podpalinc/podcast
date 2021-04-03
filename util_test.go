package podcast_test

import (
	"testing"

	podcast "github.com/podpalinc/rss-feed-generator"
	"github.com/stretchr/testify/assert"
)

func TestGenerateStringEmpty(t *testing.T) {
	t.Parallel()

	result := podcast.GenerateFeedString("")

	assert.Len(t, result, 0)
}

func TestGenerateStringAmpersand(t *testing.T) {
	t.Parallel()

	result := podcast.GenerateFeedString("Kids & Family")

	assert.Equal(t, result, "Kids &amp; Family")

}

func TestGenerateStringLessThan(t *testing.T) {
	t.Parallel()

	result := podcast.GenerateFeedString("1 < 3")

	assert.Equal(t, result, "1 &lt; 3")
}

func TestGenerateStringGreaterThan(t *testing.T) {
	t.Parallel()

	result := podcast.GenerateFeedString("3 > 1")

	assert.Equal(t, result, "3 &gt; 1")
}

func TestGenerateStringApostrophe(t *testing.T) {
	t.Parallel()

	result := podcast.GenerateFeedString("Sophie's choice")

	assert.Equal(t, result, "Sophie&apos;s choice")
}

func TestGenerateStringQuotation(t *testing.T) {
	t.Parallel()

	result := podcast.GenerateFeedString("He said \"what\"")

	assert.Equal(t, result, "He said &quot;what&quot;")
}

func TestGenerateStringCopyrightSign(t *testing.T) {
	t.Parallel()

	result := podcast.GenerateFeedString("© Podpal Inc, 2020")

	assert.Equal(t, result, "&#xA9; Podpal Inc, 2020")
}

func TestGenerateStringSoundRecordingCopyright(t *testing.T) {
	t.Parallel()

	result := podcast.GenerateFeedString("℗ Podpal Inc, 2020")

	assert.Equal(t, result, "&#x2117; Podpal Inc, 2020")
}

func TestGenerateStringTrademark(t *testing.T) {
	t.Parallel()

	result := podcast.GenerateFeedString("™ Podpal Inc, 2020")

	assert.Equal(t, result, "&#x2122; Podpal Inc, 2020")
}
