package podcast

import "encoding/xml"

/**
* This description will be used at both the channel and item level to provide the ability to add hyperlinks, formatted text, etc.
 */
type Description struct {
	XMLName xml.Name `xml:"description"`
	Text    string   `xml:",cdata"`
}
