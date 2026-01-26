package command

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cardvark/blog-aggregator/internal/database"
	"github.com/cardvark/blog-aggregator/internal/rss"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		fmt.Println("Error: feed (1) name and (2) URL required.")
		os.Exit(1)
		return fmt.Errorf("Error: no arguments found for %s:\n", cmd.name)
	}

	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	user_id := user.ID

	now := time.Now()
	nullableTime := database.GetNullTime(now)

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: nullableTime,
		Name:      feedName,
		Url:       feedURL,
		UserID:    user_id,
	}

	feed, err := s.db.CreateFeed(
		context.Background(),
		feedParams,
	)
	if err != nil {
		fmt.Printf("Error creating feed entry: %v\n", err)
		return err
	}

	err = executeFollow(s, feedURL, user)
	if err != nil {
		return err
	}

	fmt.Printf("%#v", feed)
	return nil
}

func handlerAgg(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 1 {
		fmt.Println("Error: time between reqs required.")
		os.Exit(1)
		return fmt.Errorf("Error: no arguments found for %s:\n", cmd.name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		fmt.Printf("Error parsing duration: %v\n", err)
		return err
	}
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s, user)
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(
		context.Background(),
	)
	if err != nil {
		fmt.Printf("Error retrieving feeds: %v", err)
		os.Exit(1)
		return err
	}

	for _, feed := range feeds {
		feed_user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			fmt.Printf("Error retrieving user name.")
			return err
		}
		fmt.Printf("Feed name: %s, URL: %s, user: %s\n", feed.Name, feed.Url, feed_user.Name)
	}

	return nil
}

func scrapeFeeds(s *state, user database.User) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background(), user.ID)
	if err != nil {
		fmt.Printf("Error getting next feed to fetch: %v\n", err)
		return err
	}

	now := time.Now()
	nullableTime := database.GetNullTime(now)

	feedFetchedParams := database.MarkFeedFetchedParams{
		ID:        feed.ID,
		UpdatedAt: nullableTime,
	}

	err = s.db.MarkFeedFetched(context.Background(), feedFetchedParams)
	if err != nil {
		fmt.Printf("Error marking feed as fetched: %v\n", err)
		return err
	}

	rssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	fmt.Println("Returning results from:", rssFeed.Channel.Title)
	fmt.Println("")

	for _, item := range rssFeed.Channel.Item {
		err = savePost(s, feed.ID, item)
		if err != nil {
			// return err
		}
	}

	fmt.Printf("\n")

	return nil
}
