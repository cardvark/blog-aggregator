package command

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cardvark/blog-aggregator/internal/database"
	"github.com/cardvark/blog-aggregator/internal/rss"
	"github.com/google/uuid"
)

// RSSItem{Title:"Scientists identify brain waves that define the limits of 'you'", Link:"https://www.sciencealert.com/scientists-identify-brain-waves-that-define-the-limits-of-you", Description:"<a href=\"https://news.ycombinator.com/item?id=46760099\">Comments</a>", PubDate:"Mon, 26 Jan 2026 00:10:42 +0000"}

// func saveFeed(s *state, rssFeed rss.RSSFeed, user database.User) error {
// 	for _, item := range rss.Feed.Channel.Item {
// 		savePost(s, item, rssFeed.feedID, user)
// 	}
// }

func savePost(s *state, feedID uuid.UUID, rssItem rss.RSSItem) error {

	now := time.Now()
	nullableTime := database.GetNullTime(now)

	publishedTime, err := time.Parse(time.RFC1123Z, rssItem.PubDate)
	if err != nil {
		fmt.Printf("Error parsing time: (%v) from post pubdate: %v\n", rssItem.PubDate, err)
		return err
	}

	publishedNull := database.GetNullTime(publishedTime)
	descNull := database.GetNullText(rssItem.Description)

	postParams := database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   now,
		UpdatedAt:   nullableTime,
		Title:       rssItem.Title,
		Description: descNull,
		Url:         rssItem.Link,
		PublishedAt: publishedNull,
		FeedID:      feedID,
	}

	isDupe := false

	err = s.db.CreatePost(
		context.Background(),
		postParams,
	)
	if err != nil {
		dupeEntry := "pq: duplicate key value violates unique constraint"
		if strings.Contains(fmt.Sprintf("%v", err), dupeEntry) {
			isDupe = true
		} else {
			fmt.Printf("Error creating post entry: %v\n", err)
			return err
		}
	}

	if !isDupe {
		fmt.Printf("Saved Post!\nTitle: %v\nDescription: %v\n", rssItem.Title, rssItem.Description)
	}

	return nil
}
