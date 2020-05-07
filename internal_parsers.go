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

func init() {
	MetaImageNames = []string{"og:image", "twitter:image"}
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

	if image != "" && strings.HasPrefix(image, "/") { // Image is set and is a relative URL
		image = link.Host + image // Prepend host
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
