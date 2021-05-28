package podcast

import "encoding/xml"

/**
* Represents a google play image. Is is a required tag according to https://support.google.com/podcast-publishers/answer/9889544?hl=en
 */
type GooglePlayImage struct {
	XMLName xml.Name `xml:"googleplay:image"`
	HREF    string   `xml:"href,attr"`
}

type GooglePlayCategory struct {
	XMLName              xml.Name `xml:"googleplay:category"`
	Text                 string   `xml:"text,attr"`
	GooglePlayCategories []*GooglePlayCategory
}
