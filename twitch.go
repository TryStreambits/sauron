// This file contains our Twitch parser

package sauron

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ChannelRequestJSON string
var ClipRequestJSON string

func init() {
	ChannelRequestJSON = `[{"operationName":"ChannelRoot_Channel","variables":{"currentChannelLogin":"CHANNEL","includeChanlets":true},"extensions":{"persistedQuery":{"version":1,"sha256Hash":"ce18f2832d12cabcfee42f0c72001dfa1a5ed4a84931ead7b526245994810284"}}},{"operationName":"ChannelPage_ChannelHeader","variables":{"login":"CHANNEL"},"extensions":{"persistedQuery":{"version":1,"sha256Hash":"836472cb842531bb09f1c42ef5ce40533ac215385a3b80dffcf513c8de67133a"}}}]`
	ClipRequestJSON = `[{"operationName":"ChannelRoot_Clip","variables":{"slugID":"CLIPSLUG","includeChanlets":false},"extensions":{"persistedQuery":{"version":1,"sha256Hash":"11627b974d3926baf0aaf48b85d9f122d53760b0c3e7cab1fce17b0ccb3eef2d"}}}]`
}

// Twitch is our internal Twitch parser
// This parser will leverage Twitch's GQL (used during info fetching for page content generation) to get various JSON data for the request
func Twitch(_doc *goquery.Document, url *url.URL, fullURL string) (link *Link, parserErr error) {
	link = &Link{
		Description: "",                      // Create an empty description for now
		Favicon:     "",                      // Create an empty favicon for now
		Host:        url.Host,                // Set to our provided host
		Title:       "Twitch",                // Set to standard Twitch title
		URI:         fullURL,                 // Set to provided URL
		Extras:      make(map[string]string), // Create an empty map
	}

	var content string

	if strings.Contains(url.Path, "/clip/") { // Is a clip URL
		clipSplit := strings.Split(strings.TrimSuffix(url.Path, "/"), "/") // Split on /
		link.Extras["ClipSlug"] = clipSplit[len(clipSplit)-1]              // Get the last item
		link.Extras["IsClip"] = "true"                                     // Indicate it is a clip

		content = strings.Replace(ClipRequestJSON, "CLIPSLUG", link.Extras["ClipSlug"], -1)
	} else { // Hopefully a streamer
		content = strings.Replace(ChannelRequestJSON, "CHANNEL", strings.Replace(url.Path, "/", "", -1), -1) // Replace our preset CHANNEL string
	}

	request, requestNewErr := http.NewRequest("POST", "https://gql.twitch.tv/gql", bytes.NewBuffer([]byte(content)))

	if requestNewErr != nil {
		parserErr = requestNewErr
		return
	}

	request.Header.Set("Accept-Language", RequestLanguage)            // Prefer English
	request.Header.Set("Client-Id", "kimne78kx3ncx6brgo4mv6wki5h1ko") // Generic Twitch Client-Id
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Trident/7.0; rv:11.0) like Gecko") // Fake old browser
	request.Header.Set("X-Device-Id", "derpnotreal")

	client := http.Client{
		Timeout: time.Second * 15, // 15 seconds
	}

	response, getErr := client.Do(request)

	if getErr != nil {
		parserErr = getErr
		return
	}

	defer response.Body.Close()
	responseContent, _ := ioutil.ReadAll(response.Body) // Read the body contents

	if response.StatusCode != 200 { // Status not OK
		parserErr = errors.New(string(responseContent))
		return
	}

	if link.Extras["IsClip"] == "true" { // Is a clip
		var gqlResponse []TwitchGqlClipResponse
		if parserErr = json.Unmarshal(responseContent, &gqlResponse); parserErr != nil { // Failed to parse the content as a TwitchGqlResponse
			return
		}

		if len(gqlResponse) == 0 {
			return
		}

		gql := gqlResponse[0] // Use our first response

		link.Extras["Streamer"] = gql.Data.Clip.Broadcaster.DisplayName
		link.Extras["ClipName"] = gql.Data.Clip.Title
		link.Extras["ClipSlug"] = gql.Data.Clip.Slug
		link.Title = fmt.Sprintf("%s - %s - Twitch", link.Extras["Streamer"], link.Extras["ClipName"])

		link.Extras["Game"] = gql.Data.Clip.Game.Name
		link.Extras["GameLink"] = fmt.Sprintf("https://www.twitch.tv/directory/game/%s", link.Extras["Game"])
		link.Extras["GameArtSmall"] = strings.Replace(gql.Data.Clip.Game.BoxArtURL, "-138x190", "-285x380", -1)
		link.Extras["GameArtFull"] = strings.Replace(link.Extras["GameArtSmall"], "-285x380", "", -1)
	} else { // Not a clip, treat as a video
		var gqlResponse []TwitchGqlChannelResponse                                       // Define gqlResponse as our response object
		if parserErr = json.Unmarshal(responseContent, &gqlResponse); parserErr != nil { // Failed to parse the content as a TwitchGqlResponse
			return
		}

		if len(gqlResponse) == 0 {
			return
		}

		gql := gqlResponse[0] // Use our first response

		if gql.Data.User.DisplayName == "" { // No Streamer name, which means this isn't a channel
			delete(link.Extras, "IsClip") // Delete our IsClip extras info
			return // Stick with primitive data
		}

		link.Extras["Streamer"] = gql.Data.User.DisplayName
		link.Title = fmt.Sprintf("%s - Twitch", link.Extras["Streamer"]) // Change to streamer name - Twitch
		link.Extras["IsClip"] = "false"                                  // Indicate it is not a clip

		link.Extras["StreamTitle"] = gql.Data.User.BroadcastSettings.Title
		link.Extras["Game"] = gql.Data.User.BroadcastSettings.Game.Name
		link.Extras["GameLink"] = fmt.Sprintf("https://www.twitch.tv/directory/game/%s", link.Extras["Game"])
		link.Extras["GameArtSmall"] = strings.Replace(gql.Data.User.BroadcastSettings.Game.BoxArtURL, "-85x113", "-285x380", -1)
		link.Extras["GameArtFull"] = strings.Replace(link.Extras["GameArtSmall"], "-285x380", "", -1)

		if gql.Data.Stream.Type == "live" { // Stream is active
			link.Extras["Live"] = "true"
		} else {
			link.Extras["Live"] = "false"
		}
	}

	return
}
