package setup

import "net/url"

func getTag(parsedUrl *url.URL, defaultTag string) string {
	if nil == parsedUrl {
		return ""
	}
	if parsedUrl.Query().Has("tag") {
		return parsedUrl.Query().Get("tag")
	}
	return defaultTag
}
