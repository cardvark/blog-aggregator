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

	tempCmd := command{
		args: []string{feedURL},
	}

	err = handlerFollow(s, tempCmd)
	if err != nil {
		fmt.Println("Error following new feed: ", err)
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

func handlerFeeds(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeeds(
		context.Background(),
	)
	if err != nil {
		fmt.Printf("Error retrieving feeds: %v", err)
		os.Exit(1)
		return err
	}

	for _, feed := range feeds {
		fmt.Printf("Feed name: %s, URL: %s, user: %s\n", feed.Name, feed.Url, user.Name)
	}

	return nil
}
