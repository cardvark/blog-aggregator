package command

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/cardvark/blog-aggregator/internal/database"
	"github.com/cardvark/blog-aggregator/internal/rss"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		fmt.Println("Error: URL required.")
		os.Exit(1)
		return fmt.Errorf("Error: no arguments found for %s:\n", cmd.name)
	}

	feedURL := cmd.args[0]

	feed, err := s.db.GetFeedByURL(
		context.Background(),
		feedURL,
	)
	if err != nil {
		fmt.Printf("Error retrieving feed by url: %v\n", err)
		return err
	}

	userName := s.config.Current_user_name

	user, err := s.db.GetUser(
		context.Background(),
		userName,
	)
	if err != nil {
		fmt.Println("Current user not found in database.")
		os.Exit(1)
		return err
	}

	user_id := user.ID
	feed_id := feed.ID

	now := time.Now()
	var nullableTime sql.NullTime
	nullableTime.Time = now
	nullableTime.Valid = true

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: nullableTime,
		UserID:    user_id,
		FeedID:    feed_id,
	}

	feedFollow, err := s.db.CreateFeedFollow(
		context.Background(),
		feedFollowParams,
	)
	if err != nil {
		fmt.Printf("Error inserting feedfollow: %v\n", err)
		return err
	}

	fmt.Printf("%v\n", feedFollow)
	fmt.Printf("%s is now following feed: %s.\n", userName, feedFollow.FeedName)

	return nil

}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		fmt.Println("Error: feed (1) name and (2) URL required.")
		os.Exit(1)
		return fmt.Errorf("Error: no arguments found for %s:\n", cmd.name)
	}

	feedName := cmd.args[0]
	feedURL := cmd.args[1]

	userName := s.config.Current_user_name

	user, err := s.db.GetUser(
		context.Background(),
		userName,
	)
	if err != nil {
		fmt.Println("Current user not found in database.")
		os.Exit(1)
		return err
	}

	user_id := user.ID

	now := time.Now()
	var nullableTime sql.NullTime
	nullableTime.Time = now
	nullableTime.Valid = true

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

	fmt.Printf("%#v", feed)
	return nil
}

func handlerAgg(s *state, cmd command) error {

	feedURL := "https://www.wagslane.dev/index.xml"

	rssFeed, err := rss.FetchFeed(
		context.Background(),
		feedURL,
	)
	if err != nil {
		return err
	}

	fmt.Print(rssFeed)

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
		user, err := s.db.GetUserByID(
			context.Background(),
			feed.UserID,
		)
		if err != nil {
			fmt.Printf("Error retrieving user from feed user_id: %v\n", err)
			// os.Exit(1)
			return err
		}
		fmt.Printf("Feed name: %s, URL: %s, user: %s\n", feed.Name, feed.Url, user.Name)
	}

	return nil
}
