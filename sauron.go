package sauron

import (
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/url"
	"strings"
)

// HasOverriddenInternals is a map of our internal parsers and if they have been overridden
var HasOverriddenInternals map[string]bool

// HostToParsers is our map of hostnames to custom parsers
var HostToParsers map[string]LinkParser

// RequestLanguage is the desired language to request a page with. Defaults to en-US / en
var RequestLanguage string

// UserAgent is the desired User Agent to report to a page via request. Defaults to Sauron Bot $VERSION (e.g. Sauron Bot 0.1)
var UserAgent string

const (
	// HostAlreadyRegistered is an error message for when host already has registered parser
	HostAlreadyRegistered = "Host already has a registered parser"

	// NoResponse is an error message for when we fail to get a response from a page. This may occur for timeouts.
	NoResponse = "No response from client to page"

	// PageContentNotValid is an error message for when the page requested is not HTML
	PageContentNotValid = "Page content provided is not valid HTML"

	// PageNotAccessible is an error message for when we get a non-200 status from a page
	PageNotAccessible = "Page not accessible"
)

func init() {
	HasOverriddenInternals = map[string]bool{
		"reddit.com":      false,
		"clips.twitch.tv": false,
		"twitch.tv":       false,
		"www.twitch.tv":   false,
		"youtube.com":     false,
		"youtu.be":        false,
	}

	HostToParsers = map[string]LinkParser{
		"old.reddit.com":  Reddit,
		"reddit.com":      Reddit,
		"clips.twitch.tv": Twitch,
		"twitch.tv":       Twitch,
		"www.twitch.tv":   Twitch,
		"youtu.be":        Youtube,
		"youtube.com":     Youtube,
	}

	RequestLanguage = "en-US,en;q=0.5"
	UserAgent = "Sauron Bot 0.1"
}

// ForceRegister will force register a LinkParser against the provided hostname
// This is identical to calling Unregister then Register.
func ForceRegister(hostName string, parser LinkParser) error {
	Unregister(hostName)
	return Register(hostName, parser)
}

// GetLink will get the link information for the provided url
func GetLink(urlPath string) (link *Link, parseErr error) {
	var u *url.URL              // url struct to pass to parsers
	var urlForDocument *url.URL // urlForDocument is explicitly used for document fetching.

	u, parseErr = url.Parse(urlPath) // Parse the provided URL

	if parseErr != nil { // Failed to parse the provided url
		return
	}

	if strings.HasSuffix(u.Host, "reddit.com") && !HasOverridden("reddit.com") && u.Host != "old.reddit.com" { // If the host is Reddit and our internal parser has not been overridden
		oldFriendlyURL := strings.Replace(u.String(), u.Host, "old.reddit.com", -1) // Convert host to old.reddit.com
		urlForDocument, parseErr = url.Parse(oldFriendlyURL)
	} else if u.Host == "youtu.be" && !HasOverridden("youtu.be") { // If the host is the shortened YouTube URL and our internal parser has not been overridden
		videoPath := strings.TrimPrefix(u.Path, "/")         // Trim the / from the start of the path
		urlPath = "https://youtube.com/watch?v=" + videoPath // Correct urlPath
		urlForDocument, parseErr = url.Parse(urlPath)        // Change url to more accurate struct
	} else if strings.HasSuffix(u.Host, "youtube.com") && !HasOverridden("youtube.com") && u.Host != "youtube.com" { // If the host is Youtube and our internal parser has not been overridden
		normalYoutubeURL := strings.Replace(u.String(), u.Host, "youtube.com", -1) // Convert host to youtube.com
		urlForDocument, parseErr = url.Parse(normalYoutubeURL)
	} else {
		urlForDocument, parseErr = url.Parse(u.String()) // Just duplicate u to urlForDocument
	}

	if parseErr != nil { // If we had parse errors from custom host handling
		return
	}

	client, request := NewHTTPClient(urlForDocument)
	response, getErr := client.Do(&request)

	if getErr != nil { // Failed to get a response
		parseErr = errors.New(NoResponse)
		return
	}

	if response.StatusCode != 200 && response.StatusCode != 304 { // Page is not accessible or is not unmodified
		parseErr = errors.New(PageNotAccessible)
		return
	}

	contentType := response.Header.Get("content-type")
	isHTML := strings.HasPrefix(contentType, "text/html")
	isImage := strings.HasPrefix(contentType, "image/")
	isVideo := strings.HasPrefix(contentType, "video/")

	if !isHTML && !isImage && !isVideo { // If this is not an HTML page or supported direct link
		parseErr = errors.New(PageContentNotValid)
		return
	}

	if isImage || isVideo { // If this is an image or video direct link
		extras := make(map[string]string)

		if isImage { // If this is an image
			extras["IsImageLink"] = "true"
		} else if isVideo { // If this is a video
			extras["IsVideoLink"] = "true"
		} // Intentionally use else if so we can just continue to extend it in the future

		link = &Link{
			Description: "",
			Favicon:     "",
			Host:        u.Host,
			Title:       "",
			URI:         urlPath,
			Extras:      extras,
		}
	} else if isHTML { // If this is an HTML page
		pageContent, readErr := ioutil.ReadAll(response.Body) // Read the body
		response.Body.Close()

		if readErr != nil { // If we failed to read page content
			parseErr = errors.New(PageContentNotValid)
			return
		}

		var doc *goquery.Document
		doc, parseErr = goquery.NewDocumentFromReader(bytes.NewReader(pageContent))

		if parseErr != nil { // If we failed to create a new document
			parseErr = errors.New(PageContentNotValid)
			return
		}

		if fnForDoc, fnForDocParserExists := HostToParsers[urlForDocument.Host]; fnForDocParserExists { // If we have a parser for our document
			link, parseErr = fnForDoc(doc, urlForDocument, urlPath) // Pass along to our function
			return
		} else if fnNoDoc, fnParserExists := HostToParsers[u.Host]; fnParserExists { // If we have a parser for our non-parsed / handled URL
			link, parseErr = fnNoDoc(doc, u, urlPath) // Pass along to our function
			return
		} else { // No handler
			link, parseErr = Primitive(doc, u, urlPath) // Pass along to our primitive parser
		}
	}

	return
}

// HasOverridden will check if our internal parsers have been overridden
func HasOverridden(host string) (overridden bool) {
	if overrideVal, overrideExists := HasOverriddenInternals[host]; overrideExists {
		overridden = overrideVal
	}

	return
}

// Register will attempt to register the provided parser for a specific hostname
// Hostname can be an exact match, such as "google.com" or regex.
// Attempting to register when a LinkParser is already associated will return an error.
func Register(hostName string, parser LinkParser) (regErr error) {
	if _, registered := HostToParsers[hostName]; !registered { // If this hostname has not yet been registered with a LinkParser
		HostToParsers[hostName] = parser // Add this parser

		if _, hasOverrideValue := HasOverriddenInternals[hostName]; hasOverrideValue { // Check if we have a respective entry in overridden internals for this hostname
			HasOverriddenInternals[hostName] = true
		}
	} else {
		regErr = errors.New(HostAlreadyRegistered)
	}

	return
}

// SetRequestLanguage will set the Accept-Language header for page requests
// This does not necessarily mean the page supports the language or will return with that language
func SetRequestLanguage(lang string) error {
	if lang == "" { // If the language is empty
		return errors.New("language must not be empty")
	}

	RequestLanguage = lang
	return nil
}

func SetUserAgent(agent string) error {
	if agent == "" {
		return errors.New("user agent must not be empty")
	}

	UserAgent = agent
	return nil
}

// Unregister will unregister a LinkParser with the specified hostname
func Unregister(hostName string) {
	delete(HostToParsers, hostName)
}
