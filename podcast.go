package podcast

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/pkg/errors"
)

// Constants to use while generating podcast feed.
const (
	pVersion     = "1.3.1"
	HEADER       = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
	ATOMNS       = "http://www.w3.org/2005/Atom"
	ITUNESNS     = "http://www.itunes.com/dtds/podcast-1.0.dtd"
	GOOGLEPLAYNS = "http://www.google.com/schemas/play-podcasts/1.0"
	SPOTIFYNS    = "http://www.spotify.com/ns/rss"
	CONTENT      = "http://purl.org/rss/1.0/modules/content/"
)

// Podcast represents a podcast.
type Podcast struct {
	XMLName            xml.Name `xml:"channel"`
	AtomLink           *AtomLink
	Generator          string `xml:"generator,omitempty"`
	Title              string `xml:"title"`
	Link               string `xml:"link,omitempty"`
	Description        *Description
	EncodedDescription *EncodedContent
	Language           string `xml:"language,omitempty"`
	Cloud              string `xml:"cloud,omitempty"`
	Copyright          string `xml:"copyright,omitempty"`
	Docs               string `xml:"docs,omitempty"`
	PubDate            string `xml:"pubDate,omitempty"`
	LastBuildDate      string `xml:"lastBuildDate,omitempty"`
	ManagingEditor     string `xml:"managingEditor,omitempty"`
	Rating             string `xml:"rating,omitempty"`
	SkipHours          string `xml:"skipHours,omitempty"`
	SkipDays           string `xml:"skipDays,omitempty"`
	TTL                int    `xml:"ttl,omitempty"`
	WebMaster          string `xml:"webMaster,omitempty"`
	Image              *Image
	TextInput          *TextInput

	// https://help.apple.com/itc/podcasts_connect/#/itcb54353390
	ITitle      string `xml:"itunes:title,omitempty"`
	IAuthor     string `xml:"itunes:author,omitempty"`
	ISubtitle   string `xml:"itunes:subtitle,omitempty"`
	IType       string `xml:"itunes:type,omitempty"`
	ISummary    *ISummary
	IBlock      string `xml:"itunes:block,omitempty"`
	IImage      *IImage
	IDuration   string  `xml:"itunes:duration,omitempty"`
	IExplicit   string  `xml:"itunes:explicit,omitempty"`
	IComplete   string  `xml:"itunes:complete,omitempty"`
	INewFeedURL string  `xml:"itunes:new-feed-url,omitempty"`
	IOwner      *Author // Author is formatted for itunes as-is
	ICategories []*ICategory

	// https://support.google.com/podcast-publishers/answer/9889544?hl=en
	GooglePlayAuthor      string `xml:"googleplay:author,omitempty"`
	GooglePlayDescription string `xml:"googleplay:description,omitempty"`
	GooglePlayOwner       string `xml:"googleplay:owner,omitempty"`
	GooglePlayImage       *GooglePlayImage

	Items []*Item

	encode func(w io.Writer, o interface{}) error
}

// New instantiates a Podcast with required parameters.
//
// Nil-able fields are optional but recommended as they are formatted
// to the expected proper formats.
func New(title, link string, description Description,
	pubDate, lastBuildDate *time.Time) Podcast {
	return Podcast{
		Title:       GenerateFeedString(title),
		Link:        link,
		Description: &description,
		// setup dependency (could inject later)
		encode: encoder,
	}
}

func (p *Podcast) AddTitle(title string) {
	if len(title) <= 0 {
		return
	}

	p.Title = GenerateFeedString(title)
}

// AddAuthor adds the specified Author to the podcast.
// func (p *Podcast) AddAuthor(name, email string) {
// 	if len(email) == 0 {
// 		return
// 	}
// 	p.ManagingEditor = parseAuthorNameEmail(&Author{
// 		Name:  name,
// 		Email: email,
// 	})
// 	p.IAuthor = p.ManagingEditor
// }
func (p *Podcast) AddAuthor(authors []string) {
	combinedAuthors := ""
	for i, author := range authors {
		if i == len(authors)-1 {
			combinedAuthors += author
		} else {
			combinedAuthors += author + ", "
		}
	}

	author := GenerateFeedString(combinedAuthors)
	p.IAuthor = author
	p.GooglePlayAuthor = author
}

// AddAtomLink adds a FQDN reference to an atom feed.
func (p *Podcast) AddAtomLink(href string) {
	if len(href) == 0 {
		return
	}
	p.AtomLink = &AtomLink{
		HREF: href,
		Rel:  "self",
		Type: "application/rss+xml",
	}
}

// AddCategory adds the category to the Podcast.
//
// ICategory can be listed multiple times.
//
// Calling this method multiple times will APPEND the category to the existing
// list, if any, including ICategory.
//
// Note that Apple iTunes has a specific list of categories that only can be
// used and will invalidate the feed if deviated from the list.  That list is
// as follows.
//
//   * Arts
//     * Design
//     * Fashion & Beauty
//     * Food
//     * Literature
//     * Performing Arts
//     * Visual Arts
//   * Business
//     * Business News
//     * Careers
//     * Investing
//     * Management & Marketing
//     * Shopping
//   * Comedy
//   * Education
//     * Education Technology
//     * Higher Education
//     * K-12
//     * Language Courses
//     * Training
//   * Games & Hobbies
//     * Automotive
//     * Aviation
//     * Hobbies
//     * Other Games
//     * Video Games
//   * Government & Organizations
//     * Local
//     * National
//     * Non-Profit
//     * Regional
//   * Health
//     * Alternative Health
//     * Fitness & Nutrition
//     * Self-Help
//     * Sexuality
//   * Kids & Family
//   * Music
//   * News & Politics
//   * Religion & Spirituality
//     * Buddhism
//     * Christianity
//     * Hinduism
//     * Islam
//     * Judaism
//     * Other
//     * Spirituality
//   * Science & Medicine
//     * Medicine
//     * Natural Sciences
//     * Social Sciences
//   * Society & Culture
//     * History
//     * Personal Journals
//     * Philosophy
//     * Places & Travel
//   * Sports & Recreation
//     * Amateur
//     * College & High School
//     * Outdoor
//     * Professional
//   * Technology
//     * Gadgets
//     * Podcasting
//     * Software How-To
//     * Tech News
//   * TV & Film
func (p *Podcast) AddCategory(category string, subCategories []string) {
	if len(category) == 0 {
		return
	}

	icat := ICategory{Text: category}
	for _, c := range subCategories {
		if len(c) == 0 {
			continue
		}
		icat2 := ICategory{Text: c}
		icat.ICategories = append(icat.ICategories, &icat2)
	}

	p.ICategories = append(p.ICategories, &icat)
}

func (p *Podcast) AddCopyright(copyright string) {
	if len(copyright) == 0 {
		return
	}
	p.Copyright = GenerateFeedString(copyright)
}

type podcastCategory struct {
	Name            string
	ParentCategory  string
	ChildCategories []string
}

// Utility function to parse categories and divide them into
func ParseCategories(categories []string) map[string][]string {
	PODCAST_CATEGORIES := []podcastCategory{}
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "Arts", ChildCategories: []string{"Books", "Design", "Fashion & Beauty", "Food", "Performing Arts", "Visual Arts"}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "Business", ChildCategories: []string{"Careers", "Entrepreneurship", "Investing", "Management", "Marketing", "Non-Profit"}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "Comedy", ChildCategories: []string{"Comedy Interviews", "Improv", "Stand-Up"}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "Education", ChildCategories: []string{"Courses", "How To", "Language Learning", "Self-Improvement"}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "Fiction", ChildCategories: []string{"Comedy Fiction", "Drama", "Science Fiction"}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "Government", ChildCategories: []string{}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "History", ChildCategories: []string{}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "Health & Fitness", ChildCategories: []string{"Alternative Health", "Fitness", "Medicine", "Mental Health", "Nutrition", "Sexuality"}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "Kids & Family", ChildCategories: []string{"Education for Kids", "Parenting", "Pets & Animals", "Stories for Kids"}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "Leisure", ChildCategories: []string{"Animation & Manga", "Automotive", "Aviation", "Crafts", "Games", "Hobbies", "Home & Garden", "Video Games"}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "Music", ChildCategories: []string{"Music Commentary", "Music History", "Music Interviews"}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "News", ChildCategories: []string{"Business News", "Daily News", "Entertainment News", "News Commentary", "Politics", "Sports News", "Tech News"}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "Religion & Spirituality", ChildCategories: []string{"Buddhism", "Christianity", "Hinduism", "Islam", "Judaism", "Religion", "Spirituality"}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "Science", ChildCategories: []string{"Astronomy", "Chemistry", "Earth Sciences", "Life Sciences", "Mathematics", "Natural Sciences", "Nature", "Physics", "Social Sciences"}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "Society & Culture", ChildCategories: []string{"Documentary", "Personal Journals", "Philosophy", "Places & Travel", "Relationships"}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "Sports", ChildCategories: []string{"Baseball", "Basketball", "Cricket", "Fantasy Sports", "Football", "Golf", "Hockey", "Rugby", "Running", "Soccer", "Swimming", "Tennis", "Volleyball", "Wilderness", "Wrestling"}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "Technology", ChildCategories: []string{}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "True Crime", ChildCategories: []string{}})
	PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: "TV & Film", ChildCategories: []string{"After Shows", "Film History", "Film Interviews", "Film Reviews", "TV Reviews"}})

	for _, pc := range PODCAST_CATEGORIES {
		for _, child := range pc.ChildCategories {
			PODCAST_CATEGORIES = append(PODCAST_CATEGORIES, podcastCategory{Name: child, ParentCategory: pc.Name})
		}
	}

	var findPodcastCategory = func(category string) *podcastCategory {
		for _, pc := range PODCAST_CATEGORIES {
			if &pc != nil && pc.Name == category {
				return &pc
			}
		}
		return nil
	}

	parsedCategories := make(map[string][]string)

	for _, category := range categories {
		pc := findPodcastCategory(category)
		if &pc != nil && pc != nil && pc.ParentCategory != "" {
			parentCat := GenerateFeedString(pc.ParentCategory)
			parsedCategories[parentCat] = append(parsedCategories[parentCat], strings.Replace(pc.Name, "&", "&amp;", -1))
		} else {
			cat := GenerateFeedString(pc.Name)
			parsedCategories[cat] = []string{}
		}
	}

	return parsedCategories
}

func (p *Podcast) AddDescription(description Description) {
	if len(description.Text) <= 0 {
		return
	}

	p.Description = &description
	p.EncodedDescription = &EncodedContent{
		Text: description.Text,
	}
}

func (p *Podcast) AddGenerator(generator string) {
	if len(generator) <= 0 {
		return
	}

	p.Generator = generator
}

func (p *Podcast) AddLastBuildDate(datetime string) {
	if len(datetime) == 0 {
		return
	}

	p.LastBuildDate = datetime
}

// Podcast Language Codes.
// en: "English",
// af: "Afrikaans",
// sq: "Albanian",
// an: "Aragonese",
// ar: "Arabic",
// hy: "Armenian",
// as: "Assamese",
// az: "Azerbaijani",
// eu: "Basque",
// be: "Belarusian",
// bn: "Bengali",
// bs: "Bosnian",
// br: "Breton",
// bg: "Bulgarian",
// my: "Burmese",
// ca: "Catalan",
// ch: "Chamorro",
// ce: "Chechen",
// zh: "Chinese",
// cv: "Chuvash",
// co: "Corsican",
// cr: "Cree",
// hr: "Croatian",
// cs: "Czech",
// da: "Danish",
// nl: "Dutch",
// eo: "Esperanto",
// et: "Estonian",
// fo: "Faeroese",
// fj: "Fijian",
// fi: "Finnish",
// fr: "French",
// fy: "Frisian",
// gd: "Gaelic",
// gl: "Galacian",
// ka: "Georgian",
// de: "German",
// el: "Greek",
// gu: "Gujurati",
// ht: "Haitian",
// he: "Hebrew",
// hi: "Hindi",
// hu: "Hungarian",
// is: "Icelandic",
// id: "Indonesian",
// iu: "Inuktitut",
// ga: "Irish",
// it: "Italian",
// ja: "Japanese",
// kn: "Kannada",
// ks: "Kashmiri",
// kk: "Kazakh",
// km: "Khmer",
// ky: "Kirghiz",
// ko: "Korean",
// la: "Latin",
// lv: "Latvian",
// lt: "Lithuanian",
// lb: "Luxembourgish",
// mk: "FYRO Macedonian",
// ms: "Malay",
// ml: "Malayalam",
// mt: "Maltese",
// mi: "Maori",
// mr: "Marathi",
// mo: "Moldavian",
// nv: "Navajo",
// ng: "Ndonga",
// ne: "Nepali",
// no: "Norwegian",
// oc: "Occitan",
// or: "Oriya",
// om: "Oromo",
// fa: "Persian",
// pl: "Polish",
// pt: "Portuguese",
// pa: "Punjabi",
// qu: "Quechua",
// rm: "Rhaeto-Romanic",
// ro: "Romanian",
// ru: "Russian",
// sz: "Sami (Lappish)",
// sg: "Sango",
// sa: "Sanskrit",
// sc: "Sardinian",
// sd: "Sindhi",
// si: "Singhalese",
// sr: "Serbian",
// sk: "Slovak",
// sl: "Slovenian",
// so: "Somani",
// sb: "Sorbian",
// es: "Spanish",
// sx: "Sutu",
// sw: "Swahili",
// sv: "Swedish",
// ta: "Tamil",
// tt: "Tatar",
// te: "Teluga",
// th: "Thai",
// ts: "Tsonga",
// tn: "Tswana",
// tr: "Turkish",
// tk: "Turkmen",
// uk: "Ukrainian",
// ur: "Urdu",
// ve: "Venda",
// vi: "Vietnamese",
// vo: "Volapuk",
// wa: "Walloon",
// cy: "Welsh",
// xh: "Xhosa",
// ji: "Yiddish",
// zu: "Zulu",
func (p *Podcast) AddLanguage(language string) {
	if len(language) < 2 {
		return
	}

	p.Language = GenerateFeedString(language)
}

func (p *Podcast) AddParentalAdvisory(parentalAdvisory string) {
	if parentalAdvisory == ParentalAdvisoryExplicit {
		p.IExplicit = "yes"
	} else if parentalAdvisory == ParentalAdvisoryClean {
		p.IExplicit = "no"
	}

	return
}

// AddImage adds the specified Image to the Podcast.
//
// Podcast feeds contain artwork that is a minimum size of
// 1400 x 1400 pixels and a maximum size of 3000 x 3000 pixels,
// 72 dpi, in JPEG or PNG format with appropriate file
// extensions (.jpg, .png), and in the RGB colorspace. To optimize
// images for mobile devices, Apple recommends compressing your
// image files.
func (p *Podcast) AddImage(url string) {
	if len(url) == 0 {
		return
	}
	p.Image = &Image{
		URL:   url,
		Title: p.Title,
		Link:  p.Link,
	}
	p.IImage = &IImage{HREF: url}
	p.GooglePlayImage = &GooglePlayImage{
		HREF: url,
	}
}

// AddItem adds the podcast episode.  It returns a count of Items added or any
// errors in validation that may have occurred.
//
// This method takes the "itunes overrides" approach to populating
// itunes tags according to the overrides rules in the specification.
// This not only complies completely with iTunes parsing rules; but, it also
// displays what is possible to be set on an individual episode level – if you
// wish to have more fine grain control over your content.
//
// This method imposes strict validation of the Item being added to confirm
// to Podcast and iTunes specifications.
//
// Article minimal requirements are:
//
//   * Title
//   * Description
//   * Link
//
// Audio, Video and Downloads minimal requirements are:
//
//   * Title
//   * Description
//   * Enclosure (HREF, Type and Length all required)
//
// The following fields are always overwritten (don't set them):
//
//   * GUID
//   * PubDateFormatted
//   * AuthorFormatted
//   * Enclosure.TypeFormatted
//   * Enclosure.LengthFormatted
//
// Recommendations:
//
//   * Just set the minimal fields: the rest get set for you.
//   * Always set an Enclosure.Length, to be nice to your downloaders.
//   * Follow Apple's best practices to enrich your podcasts:
//     https://help.apple.com/itc/podcasts_connect/#/itc2b3780e76
//   * For specifications of itunes tags, see:
//     https://help.apple.com/itc/podcasts_connect/#/itcb54353390
//
func (p *Podcast) AddItem(i Item) (int, error) {
	// initial guards for required fields
	if len(i.Title) == 0 {
		return len(p.Items), errors.New("Title and Description are required")
	}
	if i.Enclosure != nil {
		if len(i.Enclosure.URL) == 0 {
			return len(p.Items),
				errors.New(i.Title + ": Enclosure.URL is required")
		}
		if i.Enclosure.TypeFormatted == enclosureDefault {
			return len(p.Items),
				errors.New(i.Title + ": Enclosure.Type is required")
		}
	} else if len(i.Link) == 0 {
		return len(p.Items),
			errors.New(i.Title + ": Link is required when not using Enclosure")
	}

	// corrective actions and overrides
	//
	// i.AuthorFormatted = parseAuthorNameEmail(i.Author)
	if i.Enclosure != nil {
		if i.GUID == nil {
			i.GUID = &GUID{IsPermaLink: true, Value: i.Enclosure.URL} // yep, GUID is the Permlink URL
		}

		if i.Enclosure.Length < 0 {
			i.Enclosure.Length = 0
		}
		i.Enclosure.LengthFormatted = strconv.FormatInt(i.Enclosure.Length, 10)
		i.Enclosure.TypeFormatted = i.Enclosure.Type.String()

		// allow Link to be set for article references to Downloads,
		// otherwise set it to the enclosurer's URL.
		if len(i.Link) == 0 {
			i.Link = i.Enclosure.URL
		}
	} else {
		i.GUID = &GUID{IsPermaLink: true, Value: i.Link} // yep, GUID is the Permlink URL
	}

	// iTunes it
	//
	if len(i.IAuthor) == 0 {
		switch {
		// case i.Author != nil:
		// 	i.IAuthor = i.Author.Email
		case len(p.IAuthor) != 0:
			// i.Author = &Author{Email: p.IAuthor}
			i.IAuthor = p.IAuthor
		case len(p.ManagingEditor) != 0:
			// i.Author = &Author{Email: p.ManagingEditor}
			i.IAuthor = p.ManagingEditor
		}
	}
	if i.IImage == nil {
		if p.Image != nil {
			i.IImage = &IImage{HREF: p.Image.URL}
		}
	}

	p.Items = append(p.Items, &i)
	return len(p.Items), nil
}

func (p *Podcast) AddItunesBlock(block string) {
	if block == "hide" {
		p.IBlock = "Yes"
	} else {
		p.IBlock = "No"
	}
}

func (p *Podcast) AddItunesComplete(complete string) {
	if complete == "complete" {
		p.IComplete = "Yes"
	} else {
		p.IComplete = "No"
	}
}

func (p *Podcast) AddItunesTitle(title string) {
	if len(title) == 0 {
		return
	}

	p.ITitle = GenerateFeedString(title)
}

func (p *Podcast) AddItunesType(showType string) {
	p.IType = showType
}

func (p *Podcast) AddLink(link string) {
	if len(link) == 0 {
		return
	}

	p.Link = link
}

func (p *Podcast) AddNewFeedURL(newFeedUrl string) {
	p.INewFeedURL = newFeedUrl
}

func (p *Podcast) AddOwner(name, email string) {
	if len(name) == 0 || len(email) == 0 {
		return
	}

	p.IOwner = &Author{
		Name:  GenerateFeedString(name),
		Email: GenerateFeedString(email),
	}

	p.GooglePlayOwner = GenerateFeedString(email)
}

func (p *Podcast) AddPubDate(datetime string) {

	if len(datetime) == 0 {
		return
	}

	p.PubDate = datetime
}

// AddSubTitle adds the iTunes subtitle that is displayed with the title
// in iTunes.
//
// Note that this field should be just a few words long according to Apple.
// This method will truncate the string to 64 chars if too long with "..."
func (p *Podcast) AddSubTitle(subTitle string) {
	count := utf8.RuneCountInString(subTitle)
	if count == 0 {
		return
	}
	if count > 64 {
		s := []rune(subTitle)
		subTitle = string(s[0:61]) + "..."
	}
	p.ISubtitle = GenerateFeedString(subTitle)
}

// AddSummary adds the iTunes summary.
//
// Limit: 4000 characters
//
// Note that this field is a CDATA encoded field which allows for rich text
// such as html links: `<a href="http://www.apple.com">Apple</a>`.
func (p *Podcast) AddSummary(summary string) {
	count := utf8.RuneCountInString(summary)
	if count == 0 {
		return
	}
	if count > 4000 {
		s := []rune(summary)
		summary = string(s[0:4000])
	}
	p.ISummary = &ISummary{
		// TODO: Perform proper string programming.
		Text: summary,
	}
}

// Bytes returns an encoded []byte slice.
// func (p *Podcast) Bytes() []byte {
// 	return []byte(p.String())
// }

// Encode writes the bytes to the io.Writer stream in RSS 2.0 specification.
func (p *Podcast) Encode(w io.Writer) error {
	if _, err := w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")); err != nil {
		return errors.Wrap(err, "podcast.Encode: w.Write return error")
	}

	// atomLink := ""
	// if p.AtomLink != nil {
	// 	atomLink = "http://www.w3.org/2005/Atom"
	// }
	wrapped := PodcastWrapper{
		ITUNESNS: ITUNESNS,
		CONTENT:  CONTENT,
		// ATOMNS:   atomLink,
		Version: "2.0",
		Channel: p,
	}
	return p.encode(w, wrapped)
}

// String encodes the Podcast state to a string.
func (p *Podcast) String() string {
	b := new(bytes.Buffer)
	if err := p.Encode(b); err != nil {
		return "String: podcast.write returned the error: " + err.Error()
	}
	return b.String()
}

// // Write implements the io.Writer interface to write an RSS 2.0 stream
// // that is compliant to the RSS 2.0 specification.
// func (p *Podcast) Write(b []byte) (n int, err error) {
// 	buf := bytes.NewBuffer(b)
// 	if err := p.Encode(buf); err != nil {
// 		return 0, errors.Wrap(err, "Write: podcast.encode returned error")
// 	}
// 	return buf.Len(), nil
// }

type PodcastWrapper struct {
	XMLName      xml.Name `xml:"rss"`
	Version      string   `xml:"version,attr"`
	ATOMNS       string   `xml:"xmlns:atom,attr,omitempty"`
	ITUNESNS     string   `xml:"xmlns:itunes,attr"`
	GOOGLEPLAYNS string   `xml:"xmlns:googleplay,attr"`
	SPOTIFYNS    string   `xml:"xmlns:spotify,attr"`
	CONTENT      string   `xml:"xmlns:content,attr"`
	Channel      *Podcast
}

func NewWrapper(p *Podcast) PodcastWrapper {
	return PodcastWrapper{
		ATOMNS:       ATOMNS,
		ITUNESNS:     ITUNESNS,
		GOOGLEPLAYNS: GOOGLEPLAYNS,
		SPOTIFYNS:    SPOTIFYNS,
		CONTENT:      CONTENT,
		Version:      "2.0",
		Channel:      p,
	}
}

var encoder = func(w io.Writer, o interface{}) error {
	e := xml.NewEncoder(w)
	e.Indent("", "  ")
	if err := e.Encode(o); err != nil {
		return errors.Wrap(err, "podcast.encoder: e.Encode returned error")
	}
	return nil
}

var parseAuthorNameEmail = func(a *Author) string {
	var author string
	if a != nil {
		author = a.Email
		if len(a.Name) > 0 {
			author = fmt.Sprintf("%s (%s)", a.Email, a.Name)
		}
	}
	return author
}
