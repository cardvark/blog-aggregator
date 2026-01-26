package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
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

	for i := range rssFeed.Channel.Item {
		titleString := rssFeed.Channel.Item[i].Title
		titleString = html.UnescapeString(titleString)
		descString := rssFeed.Channel.Item[i].Description
		descString = html.UnescapeString(descString)

		titleString, err := toMarkdown(titleString)
		if err != nil {
			fmt.Printf("%v", err)
		}
		descString, err = toMarkdown(descString)
		if err != nil {
			fmt.Printf("%v", err)
		}

		rssFeed.Channel.Item[i].Title = titleString
		rssFeed.Channel.Item[i].Description = descString
	}

	return rssFeed
}

func toMarkdown(text string) (string, error) {
	markdown, err := htmltomarkdown.ConvertString(text)
	if err != nil {
		fmt.Printf("Error converting html to markdown: %v\n", err)
		return "", err
	}
	return markdown, nil
}
