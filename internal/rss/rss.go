package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	rssFeed := &RSSFeed{}

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		feedURL,
		nil,
	)
	if err != nil {
		fmt.Printf("Error generating http request: %v\n", err)
		return rssFeed, err
	}

	req.Header.Set("User-Agent", "gator")
	req.Header.Set("Accept", "application/xml")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error retrieving RSS feed: %v\n", err)
		return rssFeed, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		fmt.Println("Status error: ", res.StatusCode)
		return rssFeed, err
	}

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return rssFeed, err
	}

	err = xml.Unmarshal(dat, rssFeed)
	if err != nil {
		fmt.Printf("Error decoding response data: %v\n", err)
		return rssFeed, err
	}

	rssFeed = processHTMLEscapes(rssFeed)

	return rssFeed, nil

}

func processHTMLEscapes(rssFeed *RSSFeed) *RSSFeed {
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)

	for _, rssItem := range rssFeed.Channel.Item {
		rssItem.Title = html.UnescapeString(rssItem.Title)
		rssItem.Description = html.UnescapeString(rssItem.Description)
	}

	return rssFeed
}
