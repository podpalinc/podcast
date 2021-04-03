package podcast

import "strings"

/*
 * Ensures that the string passed in is compliant with RSS feed requirements specified in https://help.apple.com/itc/podcasts_connect/#/itc1723472cb
 */
func GenerateFeedString(str string) string {
	str = strings.Replace(str, "&", "&amp;", -1)
	str = strings.Replace(str, "<", "&lt;", -1)
	str = strings.Replace(str, ">", "&gt;", -1)
	str = strings.Replace(str, "'", "&apos;", -1)
	str = strings.Replace(str, "\"", "&quot;", -1)
	str = strings.Replace(str, "©", "&#xA9;", -1)
	str = strings.Replace(str, "℗", "&#x2117;", -1)
	str = strings.Replace(str, "™", "&#x2122;", -1)
	return str
}
