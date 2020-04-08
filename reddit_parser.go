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
