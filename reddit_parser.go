package sauron

import (
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strconv"
)

// Reddit is our internal Reddit parser
// This parser will get page information as well as Reddit post information such as dislikes, likes, and overall score
func Reddit(doc *goquery.Document, url *url.URL, fullURL string) (link *Link, parserErr error) {
	link, parserErr = Primitive(doc, url, fullURL) // First get our link information from Primitive

	link.Extras["IsRedditLink"] = "true" // Indicate it is a Reddit link

	dislikes := doc.Find(".unvoted > .dislikes").Text()
	likes := doc.Find(".unvoted > .likes").Text()
	scoreStr := doc.Find(".unvoted > .unvoted").Text()

	link.Extras["Dislikes"] = dislikes
	link.Extras["Likes"] = likes
	link.Extras["Score"] = scoreStr

	// #region Percentage Calculation

	if dislikes != "" && likes != "" && scoreStr != "" {
		var score int

		if score, parserErr = strconv.Atoi(scoreStr); parserErr != nil { // Failed to parse our score
			return
		}

		if score != 0 { // Have a score
			var downvotes int
			if downvotes, parserErr = strconv.Atoi(dislikes); parserErr != nil {
				return
			}

			var upvotes int
			if upvotes, parserErr = strconv.Atoi(likes); parserErr != nil {
				return
			}

			percentage := int((float64(downvotes) / float64(upvotes)) * 100)

			if percentage == 0 { // 100% upvote
				percentage = 100
			}

			link.Extras["Percentage"] = strconv.Itoa(percentage) // Convert our percentage to a string
		} else { // Score of 0
			link.Extras["Percentage"] = "0" // Set to 0
		}
	}

	// #endregion

	return
}
