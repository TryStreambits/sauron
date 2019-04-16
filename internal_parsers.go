package sauron

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strconv"
	"strings"
)

// This files contains our internally supported parsers

// MetaImageNames is an array of meta names commonly associated with site images
var MetaImageNames []string

// YoutubeQueriesToExtras is query info to extra metadata
var YoutubeQueriesToExtras map[string]string

func init() {
	MetaImageNames = []string{"og:image", "twitter:image"}

	YoutubeQueriesToExtras = map[string]string{
		"i":    "Index",
		"list": "Playlist",
		"t":    "Time",
		"v":    "Video",
	}
}

// Primitive is our primitive parser
// This parser will get standard page information from the most commonly supported DOM Elements
func Primitive(doc *goquery.Document, url *url.URL, fullURL string) (link *Link, parserErr error) {
	link = &Link{
		Description: "",                       // Create an empty description for now
		Favicon:     "",                       // Create an empty favicon for now
		Host:        url.Host,                 // Set to our provided host
		Title:       doc.Find("title").Text(), // Set to standard title
		URI:         fullURL,                  // Set to provided URL
		Extras:      make(map[string]string),  // Create an empty map
	}

	// #region Description Fetching

	var description string // Set description to an empty string

	if metaDescription, hasMetaDescription := doc.Find(`meta[name="description"]`).Attr("content"); hasMetaDescription { // If we found a standard meta description
		description = metaDescription
	} else { // If we did not find a standard meta description
		description = doc.Find(`meta[name="og:description"]`).AttrOr("content", "") // Attempt to get og:description and revert to empty string
	}

	link.Description = description // Update our description

	// #endregion

	// #region Favicon Fetching

	var favicon string
	var largestSize int

	doc.Find(`link[rel="icon"]`).Each(func(index int, selection *goquery.Selection) { // For each link we found
		faviconAttr := selection.AttrOr("href", "")

		if faviconAttr != "" { // Provided icon is not a string
			if !strings.HasPrefix(faviconAttr, "http") { // Is not absolute URL
				hostPrefix := url.Scheme + "://" + url.Host

				if strings.HasPrefix(faviconAttr, "/") { // At least starts with a slash
					faviconAttr = hostPrefix + faviconAttr // Just append favicon to hostPrefix
				} else {
					faviconAttr = hostPrefix + "/" + faviconAttr // Ensure there is a slash between
				}
			}
		}

		if sizeAttr, hasSizeAttr := selection.Attr("size"); hasSizeAttr { // If this icon has a size attribute
			sizeArr := strings.Split(sizeAttr, "x") // Split on x, which is common for sizes (ex. 32x32)

			if len(sizeArr) > 0 { // If we have lengths
				if iconSize, convErr := strconv.Atoi(sizeArr[0]); convErr == nil { // Convert our string to an int
					if iconSize > largestSize && faviconAttr != "" { // If this is our largest size yet and has an icon href
						favicon = faviconAttr  // Update our favicon
						largestSize = iconSize // Update our largest size
					}
				}
			}
		} else if !hasSizeAttr && favicon == "" && faviconAttr != "" { // If no size is indicated and we haven't set an icon yet
			favicon = faviconAttr
		}
	})

	link.Favicon = favicon // Update our favicon

	// #endregion

	// #region Image Parsing

	var image string // Set image to an empty string

	for _, metaImageType := range MetaImageNames {
		if metaImage, hasMetaImage := doc.Find(`meta[name="` + metaImageType + `"]`).Attr("content"); hasMetaImage { // If we found this meta image
			image = metaImage
			break
		}
	}

	if image == "" { // If we did not find an image from the metadata
		if firstImage, hasImageOnPage := doc.Find("img").Attr("src"); hasImageOnPage { // If we found an image on the page, so just the first one we find
			image = firstImage
		}
	}

	link.Image = image

	if link.Image != "" && !strings.HasPrefix(link.Image, "http") { // If a link is relative
		prefixString := url.Scheme + ":"

		if !strings.HasPrefix(link.Image, "//") { // If it doesn't start with //, which is common for Reddit or relative to the scheme being used
			prefixString += "//" // Append //
		}

		link.Image = prefixString + link.Image // Prepend
	}

	// #endregion

	return
}

// Reddit is our internal Reddit parser
// This parser will get page information as well as Reddit post information such as dislikes, likes, and overall score
func Reddit(doc *goquery.Document, url *url.URL, fullURL string) (link *Link, parserErr error) {
	link, parserErr = Primitive(doc, url, fullURL) // First get our link information from Primitive

	dislikes := doc.Find(".unvoted > .dislikes").Text()
	likes := doc.Find(".unvoted > .likes").Text()
	score := doc.Find(".unvoted > .unvoted").Text()

	link.Extras["Dislikes"] = dislikes
	link.Extras["Likes"] = likes
	link.Extras["Score"] = score

	// #region Percentage Calculation

	if dislikes != "" && likes != "" && score != "" {
		var convertScoreErr error
		var downvotes int
		var upvotes int

		downvotes, convertScoreErr = strconv.Atoi(dislikes)

		if convertScoreErr == nil { // No error converting downvotes
			upvotes, convertScoreErr = strconv.Atoi(likes)
		}

		if convertScoreErr == nil { // No error converting downvotes and upvotes
			percentage := int((float64(downvotes) / float64(upvotes)) * 100)

			if percentage == 0 { // 100% upvote
				percentage = 100
			}

			link.Extras["Percentage"] = strconv.Itoa(percentage) // Convert our percentage to a string
		}
	}

	// #endregion

	return
}

// Youtube is our internal Youtube parser
// This parser will get page information as well as add extra metadata for various shorteners and form factors
func Youtube(doc *goquery.Document, url *url.URL, fullURL string) (link *Link, parserErr error) {
	link, parserErr = Primitive(doc, url, fullURL)            // First get our link information from Primitive
	link.Title = strings.TrimSuffix(link.Title, " - YouTube") // Strip - Youtube from the Title

	if len(url.RawQuery) != 0 { // If we have query information
		queries := strings.Split(url.RawQuery, "&") // Split on &

		for _, query := range queries { // For each query
			queryInfo := strings.Split(query, "=") // Split individual query into type and value

			if extrasType, queryTypeExists := YoutubeQueriesToExtras[queryInfo[0]]; queryTypeExists { // If this query type exists
				link.Extras[extrasType] = queryInfo[1] // Set in our extras the type to the value from queryInfo
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
