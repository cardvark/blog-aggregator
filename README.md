Simple blog aggregator program.

Requires go (v1.25+) and postgresql (v15+) installed.

Install with `go install github.com/cardvark/blog-aggregator`

Commands: blog-aggregator [command] [args]
- register [name] --> Register and log user in.
- login [name] --> Switch to another previously registered user. Must have been registered.
- users --> List all registered users.
- addfeed [name] [url] --> Add a new rss feed. Must include name of the feed and url.
- feeds --> List all feeds that have been added.
- follow [feed_url] --> For the currently logged in user, follow a feed (that has already been added), by the feed's url.
- following --> List all feeds that current user is following.
- unfollow [feed_url] --> Unfollow a given feed
- agg [time] --> Run the blog aggregator, retrieving the last updated feed every X time units. (include both value and increment. E.g. 60s, 2m, etc.)
- browse [count] --> Retrieve the last X posts that were retrieved, ordered by publishing date desc. If no argument is provided, the last 2 will be returned.
- reset --> Delete all users, feeds, feed follows, posts.
