package sauron

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
)

// This file contains our various structs for Sauron

// LinkParser is a function which takes in a parsed document, URL struct and a string, and returns a pointer to a Link or an error
type LinkParser func(*goquery.Document, *url.URL, string) (*Link, error)

// Link is our structured information about a URL provided to Sauron's Parser
type Link struct {
	Description, Favicon, Host, Title, URI string

	// Extras is our extra metadata.
	// This may be used by internal and external parsers to communicate additional information about the URL in question
	Extras map[string]string
}
