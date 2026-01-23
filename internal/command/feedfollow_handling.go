package command

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/cardvark/blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerUnfollow(s *state, cmd command, user database.User) error {
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
	removeFollowParams := database.RemoveFeedFollowForUserParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.RemoveFeedFollowForUser(
		context.Background(),
		removeFollowParams,
	)
	if err != nil {
		fmt.Printf("Error removing feed follow from db: %v\n", err)
		return err
	}

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(
		context.Background(),
		user.ID,
	)
	if err != nil {
		fmt.Printf("Error retrieving feed follows for user: %v\n", err)
		return err
	}
	fmt.Println(feedFollows)

	fmt.Printf("Current user %s is following:\n", user.Name)
	for _, feedFollow := range feedFollows {
		fmt.Printf("- %v\n", feedFollow.FeedName)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
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
	fmt.Printf("%s is now following feed: %s.\n", user.Name, feedFollow.FeedName)

	return nil
}
