package podcast

import "encoding/xml"

// EncodedContent encapsulates the recommended way to add HTML content in the description
// that is properly formatted across all podcast distributors (<content:encoded>)
type EncodedContent struct {
	XMLName xml.Name `xml:"content:encoded"`
	Text    string   `xml:",cdata"`
}
