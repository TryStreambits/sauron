package main

import (
	"fmt"
	"github.com/JoshStrobl/sauron"
	"github.com/JoshStrobl/trunk"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

func main() {
	image, imageLinkErr := sauron.GetLink("https://i3.ytimg.com/vi/OE-Y-PotqTQ/maxresdefault.jpg")

	if imageLinkErr == nil { // Got the Image
		if image.Extras["IsImageLink"] == "true" { // If we have an Extras metadata set
			trunk.LogSuccess("Got JPG")
			fmt.Printf("%v\n", image)
		} else {
			trunk.LogErr("Failed to determine that our image is precisely that.")
		}
	} else {
		trunk.LogErr(fmt.Sprintf("Failed to get Image: %v", imageLinkErr))
	}

	video, videoLinkErr := sauron.GetLink("http://mirrors.standaloneinstaller.com/video-sample/small.mp4")

	if videoLinkErr == nil { // Got the video
		if video.Extras["IsVideoLink"] == "true" { // If we have an Extras metadata set
			trunk.LogSuccess("Got MP4")
			fmt.Printf("%v\n", video)
		} else {
			trunk.LogErr("Failed to determine that our video is prcisely that.")
		}
	} else {
		trunk.LogErr(fmt.Sprintf("Failed to get Video: %v", videoLinkErr))
	}

	twitter, twitterLinkErr := sauron.GetLink("https://twitter.com/trystreambits/status/1246090584714027010")

	if twitterLinkErr == nil { // Got Twitter link data
		trunk.LogSuccess("Got @trystreambits Tweet")
		fmt.Printf("%v\n", twitter)
	} else {
		trunk.LogErr(fmt.Sprintf("Failed to get Tweet: %v", twitterLinkErr))
	}
	bigBuckBunnyLink, linkErr := sauron.GetLink("https://www.youtube.com/watch?v=YE7VzlLtp-4")

	if linkErr == nil { // Successfully got link data
		if bigBuckBunnyLink.Title == "Big Buck Bunny" && bigBuckBunnyLink.Extras["IsVideo"] == "true" { // Successfully fetched
			trunk.LogSuccess(fmt.Sprintf("Fetched Big Buck Bunny. Has the following content: %s\n", bigBuckBunnyLink))
		} else { // Details do not match
			trunk.LogErr(fmt.Sprintf("Successfully fetched Big Buck Bunny but content does not match expectation: %s\n", bigBuckBunnyLink))
		}
	} else { // If we failed to fetch Big Buck Bunny
		trunk.LogErr(fmt.Sprintf("Failed to get Big Buck Bunny: %v", linkErr))
	}

	playlistTestLink, playlistTestLinkErr := sauron.GetLink("https://www.youtube.com/playlist?list=PLFF5D72E24079FB50")

	if playlistTestLinkErr == nil { // Successfully got playlist
		if playlistTestLink.Title == "Mat Kearney - Young Love" && // Name matches
			playlistTestLink.Extras["IsPlaylist"] == "true" && // Is a Playlist
			playlistTestLink.Image == "https://i.ytimg.com/vi/FANROVxej50/hqdefault.jpg" { // Playlist Image matches
			trunk.LogSuccess(fmt.Sprintf("Fetched Youtube Playlist. Has the following content: %s\n", playlistTestLink))
		} else {
			trunk.LogErr(fmt.Sprintf("Successfully fetched Youtube Playlist but content does not match expectation: %s\n", playlistTestLink))
		}
	} else {
		trunk.LogErr(fmt.Sprintf("Failed to get Youtube Playlist: %v", playlistTestLink))
	}

	redditPost, redditLinkErr := sauron.GetLink("https://www.reddit.com/r/SolusProject/comments/b2a8x0/solus_4_fortitude_released_solus/")

	if redditLinkErr == nil { // Successfully got reddit post
		if redditPost.Title == "Solus 4 Fortitude Released | Solus : SolusProject" && redditPost.Extras["Likes"] != "" { // Successfully got Reddit post
			trunk.LogSuccess(fmt.Sprintf("Fetched Reddit post. Has the following content: %s\n", redditPost))
		} else { // Failed to get reddit post, potentially likes
			trunk.LogErr(fmt.Sprintf("Successfully fetched Reddit post but content does not match expectations: %s\n", redditPost))
		}
	} else { // Failed to fetch Reddit post
		trunk.LogErr(fmt.Sprintf("Failed to get Reddit post: %v", redditLinkErr))
	}

	sauron.Register("joshuastrobl.com", PersonalSiteHandler)

	personalSiteLink, personalLinkErr := sauron.GetLink("https://joshuastrobl.com")

	if personalLinkErr == nil { // Successfully got personal site
		if personalSiteLink.Title == "Home | Joshua Strobl" && strings.HasPrefix(personalSiteLink.Extras["Generator"], "Hugo") { // Successfully got Personal Site
			trunk.LogSuccess(fmt.Sprintf("Fetched Personal Site. Has the following content: %s\n", personalSiteLink))
		} else { // Failed to get personal site, potentially generator info
			trunk.LogErr(fmt.Sprintf("Successfully fetched Personal Site but content does not match expecations: %s\n", personalSiteLink))
		}
	} else { // Failed to get personal site
		trunk.LogErr(fmt.Sprintf("Failed to get Personal Site: %v", personalLinkErr))
	}
}

// PersonalSiteHandler is a LinkParser for joshuastrobl.com
func PersonalSiteHandler(doc *goquery.Document, u *url.URL, fullPath string) (link *sauron.Link, parseErr error) {
	link, parseErr = sauron.Primitive(doc, u, fullPath) // Handle with Primitive first

	if parseErr != nil { // If we failed with primitive
		return
	}

	generatorElem := doc.Find(`meta[name="generator"]`) // Get the meta generator tag
	generator := generatorElem.AttrOr("content", "ERROR")
	link.Extras["Generator"] = generator

	return
}
