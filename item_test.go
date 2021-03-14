package podcast_test

import (
	"testing"

	podcast "github.com/podpalinc/rss-feed-generator"
	"github.com/stretchr/testify/assert"
)

func TestAddGUIDEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{}

	i.AddGUID("")

	assert.Len(t, i.GUID, 0)
}

func TestAddGUID(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{}

	i.AddGUID("987654321e7a8183nrknmfd9asfdsg")

	assert.Equal(t, i.GUID, "987654321e7a8183nrknmfd9asfdsg")
}

func TestAddTitleEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{}

	i.AddTitle("")

	assert.Len(t, i.Title, 0)
}

func TestAddTitle(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{}

	i.AddTitle("Title")

	assert.Equal(t, i.Title, "Title")
}

func TestAddLinkEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{}

	i.AddLink("")

	assert.Len(t, i.Link, 0)
}

func TestAddLink(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{}

	link := "https://google.com"

	i.AddLink(link)

	assert.Equal(t, i.Link, link)
}

func TestAddDescriptionEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{}

	i.AddDescription("")

	assert.Len(t, i.Description, 0)
}

func TestAddDescription(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{}

	desc := "This is a sample description"

	i.AddDescription(desc)

	assert.Equal(t, i.Description, desc)
}

func TestAddEpisodeNumberInvalid(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}

	i.AddEpisodeNumber(0)

	assert.Len(t, i.EpisodeNumber, 0)
}

func TestAddEpisodeNumber(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}

	i.AddEpisodeNumber(3)

	assert.Equal(t, "3", i.EpisodeNumber)
}

func TestAddEpisodeTypeEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}

	i.AddEpisodeType("")

	assert.Len(t, i.EpisodeType, 0)
}

func TestAddEpisodeTypeInvalid(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}

	i.AddEpisodeType("invalid")

	assert.Len(t, i.EpisodeType, 0)
}

func TestAddEpisodeType(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}

	i.AddEpisodeType("full")

	assert.Equal(t, "full", i.EpisodeType)
}

func TestAddEpisodeParentalAdvisoryEmpty(t *testing.T) {
	t.Parallel()

	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}

	i.AddParentalAdvisory("invalid")

	assert.EqualValues(t, i.IExplicit, "")
}

func TestAddEpisodeParentalAdvisoryExplicit(t *testing.T) {
	t.Parallel()

	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}

	i.AddParentalAdvisory(podcast.ParentalAdvisoryExplicit)

	assert.EqualValues(t, i.IExplicit, "true")
}

func TestAddEpisodeParentalAdvisoryClean(t *testing.T) {
	t.Parallel()

	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}

	i.AddParentalAdvisory(podcast.ParentalAdvisoryClean)

	assert.EqualValues(t, i.IExplicit, "false")
}

func TestAddPubDateEmpty(t *testing.T) {
	t.Parallel()

	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}

	i.AddPubDate("")

	assert.Len(t, i.PubDate, 0)
}

func TestAddPubDate(t *testing.T) {
	t.Parallel()

	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}

	i.AddPubDate("Sun, 14 Mar 2021 18:34:05 +0000")

	assert.Equal(t, i.PubDate, "Sun, 14 Mar 2021 18:34:05 +0000")
}

func TestItemAddSummaryTooLong(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}
	summary := ""
	for {
		if len(summary) >= 4051 {
			break
		}
		summary += "abc ss 5 "
	}

	// act
	i.AddSummary(summary)

	// assert
	assert.Len(t, i.ISummary.Text, 4000)
}

func TestAddEpisodeBlockEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}
	summary := ""
	for {
		if len(summary) >= 4051 {
			break
		}
		summary += "abc ss 5 "
	}

	i.AddItunesBlock("")

	assert.Equal(t, i.IBlock, "No")
}

func TestAddEpisodeBlockHide(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}
	summary := ""
	for {
		if len(summary) >= 4051 {
			break
		}
		summary += "abc ss 5 "
	}

	i.AddItunesBlock("hide")

	assert.Equal(t, i.IBlock, "Yes")
}

func TestAddSeasonNumberInvalid(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}

	i.AddSeasonNumber(0)

	assert.Len(t, i.SeasonNumber, 0)
}

func TestAddSeasonNumber(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}

	i.AddSeasonNumber(3)

	assert.Equal(t, "3", i.SeasonNumber)
}

func TestItemAddImageEmptyUrl(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}

	// act
	i.AddImage("")

	// assert
	assert.Nil(t, i.IImage)
}

func TestItemAddDurationZero(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}
	d := int64(0)

	// act
	i.AddDuration(d)

	// assert
	assert.EqualValues(t, "", i.IDuration)
}

func TestItemAddDurationLessThanZero(t *testing.T) {
	t.Parallel()

	// arrange
	i := podcast.Item{
		Title:       "item.title",
		Description: "item.desc",
		Link:        "http://example.com/article.html",
	}
	d := int64(-13)

	// act
	i.AddDuration(d)

	// assert
	assert.EqualValues(t, "", i.IDuration)
}
