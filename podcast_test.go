package podcast_test

import (
	"errors"
	"testing"
	"time"

	podcast "github.com/podpalinc/rss-feed-generator"
	"github.com/stretchr/testify/assert"
)

var (
	createdDate = time.Date(2017, time.February, 1, 8, 21, 52, 0, time.UTC)
	updatedDate = createdDate.AddDate(0, 0, 5)
	pubDate     = createdDate.AddDate(0, 0, 3)
)

func TestNewNonNils(t *testing.T) {
	t.Parallel()

	// arrange
	ti, l, d := "title", "link", "description"

	// act
	p := podcast.New(ti, l, d, &createdDate, &updatedDate)

	// assert
	assert.EqualValues(t, ti, p.Title)
	assert.EqualValues(t, l, p.Link)
	assert.EqualValues(t, d, p.Description)
}

func TestNewNils(t *testing.T) {
	t.Parallel()

	// arrange
	ti, l, d := "title", "link", "description"

	// act
	p := podcast.New(ti, l, d, nil, nil)

	// assert
	assert.EqualValues(t, ti, p.Title)
	assert.EqualValues(t, l, p.Link)
	assert.EqualValues(t, d, p.Description)
}

func TestAddAuthorEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	// act
	p.AddAuthor([]string{})

	// assert
	// assert.Len(t, p.ManagingEditor, 0)
	assert.Len(t, p.IAuthor, 0)
}

func TestAddAuthor(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	// act
	p.AddAuthor([]string{"Joe", "Lisa", "Aaron's Woods"})

	// assert
	// assert.Len(t, p.ManagingEditor, 0)
	assert.Equal(t, p.IAuthor, "Joe, Lisa, Aaron&apos;s Woods")
}

func TestAddCopyright(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	// act
	p.AddCopyright("inserted © text")

	// assert
	assert.Equal(t, p.Copyright, "inserted &#xA9; text")
}

func TestAddCopyrightEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	// act
	p.AddCopyright("")

	// assert
	assert.Equal(t, p.Copyright, "")
}

func TestAddAtomLinkHrefEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	// act
	p.AddAtomLink("")

	// assert
	assert.Nil(t, p.AtomLink)
}

func TestAddCategoryEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	// act
	p.AddCategory("", nil)

	// assert
	assert.Len(t, p.ICategories, 0)
	assert.Len(t, p.Category, 0)
}
func TestAddLanguageEmpty(t *testing.T) {
	t.Parallel()

	p := podcast.New("title", "link", "description", nil, nil)

	p.AddLanguage("")

	assert.Len(t, p.Language, 0)
}

func TestAddCategorySubCatEmpty1(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	// act
	p.AddCategory("mycat", []string{""})

	// assert
	assert.Len(t, p.ICategories, 1)
	assert.EqualValues(t, p.Category, "mycat")
	assert.Len(t, p.ICategories[0].ICategories, 0)
}

func TestAddCategorySubCatEmpty2(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	// act
	p.AddCategory("mycat", []string{"xyz", "", "abc"})

	// assert
	assert.Len(t, p.ICategories, 1)
	assert.EqualValues(t, p.Category, "mycat")
	assert.Len(t, p.ICategories[0].ICategories, 2)
}

func TestParseCategories(t *testing.T) {
	t.Parallel()

	out := podcast.ParseCategories([]string{"Arts", "Books", "Religion & Spirituality", "Christianity", "Buddhism", "Sports", "Health & Fitness", "Documentary"})

	expected := map[string][]string{
		"Arts":                        []string{"Books"},
		"Religion &amp; Spirituality": []string{"Christianity", "Buddhism"},
		"Sports":                      []string{},
		"Health &amp; Fitness":        []string{},
		"Society &amp; Culture":       []string{"Documentary"},
	}

	assert.EqualValues(t, expected, out)
}

func TestParseCategoriesChildCatEmpty(t *testing.T) {
	t.Parallel()

	out := podcast.ParseCategories([]string{"Health & Fitness"})

	expected := map[string][]string{
		"Health &amp; Fitness": []string{},
	}

	assert.EqualValues(t, expected, out)
}

func TestAddParentalAdvisoryEmpty(t *testing.T) {
	t.Parallel()

	p := podcast.New("title", "link", "description", nil, nil)

	p.AddParentalAdvisory("invalid")

	assert.EqualValues(t, p.IExplicit, "")
}

func TestAddParentalAdvisoryExplicit(t *testing.T) {
	t.Parallel()

	p := podcast.New("title", "link", "description", nil, nil)

	p.AddParentalAdvisory(podcast.ParentalAdvisoryExplicit)

	assert.EqualValues(t, p.IExplicit, "true")
}

func TestAddParentalAdvisoryClean(t *testing.T) {
	t.Parallel()

	p := podcast.New("title", "link", "description", nil, nil)

	p.AddParentalAdvisory(podcast.ParentalAdvisoryClean)

	assert.EqualValues(t, p.IExplicit, "false")
}

func TestAddImageEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	// act
	p.AddImage("")

	// assert
	assert.Nil(t, p.Image)
}

func TestAddImage(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	// act
	p.AddImage("https://google.com")

	assert.Equal(t, p.Image.URL, "https://google.com")
}

func TestAddItemEmptyTitleDescription(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)
	i := podcast.Item{}

	// act
	added, err := p.AddItem(i)

	// assert
	assert.EqualValues(t, 0, added)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Title")
	assert.Contains(t, err.Error(), "Description")
	assert.Contains(t, err.Error(), "required")
}

func TestAddItemEmptyEnclosureURL(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)
	i := podcast.Item{Title: "title", Description: "desc"}
	i.AddEnclosure("", podcast.MP3, "audio/mpeg", 1)

	// act
	added, err := p.AddItem(i)

	// assert
	assert.EqualValues(t, 0, added)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Enclosure.URL is required")
}

func TestAddItemEmptyEnclosureType(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)
	i := podcast.Item{Title: "title", Description: "desc"}
	i.AddEnclosure("http://example.com/1.mp3", podcast.MP3, "application/octet-stream", 1)

	// act
	added, err := p.AddItem(i)

	// assert
	assert.EqualValues(t, 0, added)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Enclosure.Type is required")
}

func TestAddItemEmptyLink(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)
	i := podcast.Item{Title: "title", Description: "desc"}

	// act
	added, err := p.AddItem(i)

	// assert
	assert.EqualValues(t, 0, added)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Link is required")
}

func TestAddItemEnclosureLengthMin(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)
	i := podcast.Item{Title: "title", Description: "desc"}
	i.AddEnclosure("http://example.com/1.mp3", podcast.MP3, "audio/mpeg", -1)

	// act
	added, err := p.AddItem(i)

	// assert
	assert.EqualValues(t, 1, added)
	assert.NoError(t, err)
	assert.Len(t, p.Items, 1)
	assert.EqualValues(t, 0, p.Items[0].Enclosure.Length)
}

func TestAddItemEnclosureNoLinkOverride(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)
	i := podcast.Item{Title: "title", Description: "desc"}
	i.AddEnclosure("http://example.com/1.mp3", podcast.MP3, "audio/mpeg", -1)

	// act
	added, err := p.AddItem(i)

	// assert
	assert.EqualValues(t, 1, added)
	assert.NoError(t, err)
	assert.Len(t, p.Items, 1)
	assert.EqualValues(t, i.Enclosure.URL, p.Items[0].Link)
}

func TestAddItemEnclosureLinkPresentNoOverride(t *testing.T) {
	t.Parallel()

	// arrange
	theLink := "http://someotherurl.com/story.html"
	p := podcast.New("title", "link", "description", nil, nil)
	i := podcast.Item{Title: "title", Description: "desc"}
	i.Link = theLink
	i.AddEnclosure("http://example.com/1.mp3", podcast.MP3, "audio/mpeg", -1)

	// act
	added, err := p.AddItem(i)

	// assert
	assert.EqualValues(t, 1, added)
	assert.NoError(t, err)
	assert.Len(t, p.Items, 1)
	assert.EqualValues(t, theLink, p.Items[0].Link)
}

func TestAddItemNoEnclosureGUIDValid(t *testing.T) {
	t.Parallel()

	// arrange
	theLink := "http://someotherurl.com/story.html"
	p := podcast.New("title", "link", "description", nil, nil)
	i := podcast.Item{Title: "title", Description: "desc"}
	i.Link = theLink

	// act
	added, err := p.AddItem(i)

	// assert
	assert.EqualValues(t, 1, added)
	assert.NoError(t, err)
	assert.Len(t, p.Items, 1)
	assert.EqualValues(t, theLink, p.Items[0].GUID)
}

func TestAddItemWithEnclosureGUIDSet(t *testing.T) {
	t.Parallel()

	// arrange
	theLink := "http://someotherurl.com/story.html"
	theGUID := "someGUID"
	length := 3
	p := podcast.New("title", "link", "description", nil, nil)
	i := podcast.Item{
		Title:       "title",
		Description: "desc",
		GUID:        theGUID,
	}
	i.AddEnclosure(theLink, podcast.MP3, "audio/mpeg", int64(length))

	// act
	added, err := p.AddItem(i)

	// assert
	assert.EqualValues(t, 1, added)
	assert.NoError(t, err)
	assert.Len(t, p.Items, 1)
	assert.EqualValues(t, theGUID, p.Items[0].GUID)
	assert.EqualValues(t, length, p.Items[0].Enclosure.Length)
}

func TestAddItemAuthor(t *testing.T) {
	t.Parallel()

	// arrange
	// theAuthor := podcast.Author{Name: "Jane Doe", Email: "me@janedoe.com"}
	p := podcast.New("title", "link", "description", nil, nil)
	i := podcast.Item{Title: "title", Description: "desc", Link: "http://a.co/"}
	// i.Author = &theAuthor

	// act
	added, err := p.AddItem(i)

	// assert
	assert.EqualValues(t, 1, added)
	assert.NoError(t, err)
	assert.Len(t, p.Items, 1)
	// assert.EqualValues(t, &theAuthor, p.Items[0].Author)
	// assert.EqualValues(t, theAuthor.Email, p.Items[0].IAuthor)
}

func TestAddItemRootManagingEditorSetsAuthorIAuthor(t *testing.T) {
	t.Parallel()

	// arrange
	theAuthor := "me@janedoe.com"
	p := podcast.New("title", "link", "description", nil, nil)
	p.ManagingEditor = theAuthor
	i := podcast.Item{Title: "title", Description: "desc", Link: "http://a.co/"}

	// act
	added, err := p.AddItem(i)

	// assert
	assert.EqualValues(t, 1, added)
	assert.NoError(t, err)
	assert.Len(t, p.Items, 1)
	// assert.EqualValues(t, theAuthor, p.Items[0].Author.Email)
	assert.EqualValues(t, theAuthor, p.Items[0].IAuthor)
}

func TestAddItemRootIAuthorSetsAuthorIAuthor(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)
	p.IAuthor = "me@janedoe.com"
	i := podcast.Item{Title: "title", Description: "desc", Link: "http://a.co/"}

	// act
	added, err := p.AddItem(i)

	// assert
	assert.EqualValues(t, 1, added)
	assert.NoError(t, err)
	assert.Len(t, p.Items, 1)
	// assert.EqualValues(t, "me@janedoe.com", p.Items[0].Author.Email)
	assert.EqualValues(t, "me@janedoe.com", p.Items[0].IAuthor)
}

func TestAddBlockEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	p.AddItunesBlock("")

	assert.Equal(t, p.IBlock, "No")
}

func TestAddBlockHide(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	p.AddItunesBlock("hide")

	assert.Equal(t, p.IBlock, "Yes")
}

func TestAddCompleteEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	p.AddItunesComplete("")

	assert.Equal(t, p.IComplete, "No")
}

func TestAddComplete(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	p.AddItunesComplete("complete")

	assert.Equal(t, p.IComplete, "Yes")
}

func TestAddItunesImageEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	p.AddItunesImage("")

	assert.Nil(t, p.IImage)
}

func TestAddItunesImage(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	p.AddItunesImage("https://google.com")

	assert.Equal(t, p.IImage.HREF, "https://google.com")
}

func TestAddShowTypeEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	p.AddItunesType("")

	assert.Len(t, p.IType, 0)
}

func TestAddShowType(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	p.AddItunesType("episodic")

	assert.Equal(t, "episodic", p.IType)
}

func TestAddPodcastLinkEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("", "", "", nil, nil)

	p.AddLink("")

	assert.Len(t, p.Link, 0)
	assert.Equal(t, p.Link, "")
}

func TestAddPodcastLink(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("", "", "", nil, nil)

	p.AddLink("https://google.com")

	assert.Equal(t, p.Link, "https://google.com")
}

func TestAddNewFeedURLEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	p.AddNewFeedURL("")

	assert.Equal(t, p.INewFeedURL, "")
}

func TestAddNewFeedURL(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	p.AddNewFeedURL("https://podpal.com")

	assert.Equal(t, p.INewFeedURL, "https://podpal.com")
}
func TestAddOwnerEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	p.AddOwner("", "")

	assert.Nil(t, p.IOwner)
}

func TestAddOwner(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "link", "description", nil, nil)

	p.AddOwner("joe", "joe@podpal.com")

	assert.Equal(t, p.IOwner.Name, "joe")
	assert.Equal(t, p.IOwner.Email, "joe@podpal.com")
}

func TestAddSubTitleEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "desc", "Link", nil, nil)

	// act
	p.AddSubTitle("")

	// assert
	assert.Len(t, p.ISubtitle, 0)
}

func TestAddSubTitleTooLong(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "desc", "Link", nil, nil)
	subTitle := ""
	for {
		if len(subTitle) >= 80 {
			break
		}
		subTitle += "ajd 2 "
	}

	// act
	p.AddSubTitle(subTitle)

	// assert
	assert.Len(t, p.ISubtitle, 64)
}

func TestAddSummaryTooLong(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New(
		"title",
		"desc",
		"Link",
		nil, nil)
	summary := ""
	for {
		if len(summary) >= 4051 {
			break
		}
		summary += "jax ss 7 "
	}

	// act
	p.AddSummary(summary)

	// assert
	assert.Len(t, p.ISummary.Text, 4000)
}

func TestAddSummaryEmpty(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "desc", "Link", nil, nil)

	// act
	p.AddSummary("")

	// assert
	assert.Nil(t, p.ISummary)
}

type errWriter struct{}

func (w errWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("it was bad")
}

func TestAddPodcastTitleEmpty(t *testing.T) {
	t.Parallel()

	p := podcast.Podcast{}
	p.AddTitle("")

	assert.Len(t, p.Title, 0)
}

func TestAddPodcastTitle(t *testing.T) {
	t.Parallel()

	p := podcast.Podcast{}
	p.AddTitle("Podpal™, Inc. 20")

	assert.Equal(t, p.Title, "Podpal&#x2122;, Inc. 20")
}

func TestEncodeWriterError(t *testing.T) {
	t.Parallel()

	// arrange
	p := podcast.New("title", "desc", "Link", nil, nil)

	// act
	err := p.Encode(&errWriter{})

	// assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "w.Write return error")
}
