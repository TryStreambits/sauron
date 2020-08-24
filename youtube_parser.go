package sauron

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

// YoutubeQueriesToExtras is query info to extra metadata
var YoutubeQueriesToExtras map[string]string

func init() {
	YoutubeQueriesToExtras = map[string]string{
		"i":    "Index",
		"list": "Playlist",
		"t":    "Time",
		"v":    "Video",
	}
}

// Youtube is our internal Youtube parser
// This parser will get page information as well as add extra metadata for various shorteners and form factors
func Youtube(doc *goquery.Document, url *url.URL, fullURL string) (link *Link, parserErr error) {
	link, parserErr = Primitive(doc, url, fullURL)            // First get our link information from Primitive
	link.Title = strings.TrimSuffix(link.Title, " - YouTube") // Strip - Youtube from the Title

	if link.Title == "" { // If this has no title
		link.Title = doc.Find("meta[itemprop=\"name\"]").AttrOr("content", "YouTube")
	}

	link.Extras["IsYouTubeLink"] = "true" // Indicate it is a YouTube link

	if len(url.RawQuery) != 0 { // If we have query information
		for queryParam := range url.Query() { // For each map of query params to values
			queryVal := url.Query().Get(queryParam) // Get the first value

			if extrasType, queryTypeExists := YoutubeQueriesToExtras[queryParam]; queryTypeExists { // If this query param exists
				link.Extras[extrasType] = queryVal
			}
		}

		link.Extras["IsPlaylist"] = "false"
		link.Extras["IsVideo"] = "false"

		if strings.HasPrefix(url.Path, "/playlist") { // Is a Playlist
			link.Extras["IsPlaylist"] = "true"

			if imageURL, parseErr := url.Parse(link.Image); parseErr == nil { // Parse our link image
				imageURL.RawQuery = ""         // Clear out query
				link.Image = imageURL.String() // Convert back to string
			} else {
				parserErr = parseErr
			}
		}

		if strings.HasPrefix(url.Path, "/watch") { // Is a Video
			link.Image = "https://img.youtube.com/vi/" + link.Extras["Video"] + "/maxresdefault.jpg"
			link.Extras["IsVideo"] = "true"
		}
	}

	return
}
