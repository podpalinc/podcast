package podcast

import "encoding/xml"

// EncodedContent encapsulates the recommended way to add HTML content in the description
// that is properly formatted across all podcast distributors (<content:encoded>)
type GUID struct {
	XMLName     xml.Name `xml:"guid"`
	IsPermaLink bool     `xml:"isPermaLink,attr"`
	Value       string
}
