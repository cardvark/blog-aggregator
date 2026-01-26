
-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *
;

-- name: GetFeedByURL :one
select * 
from feeds 
where url = $1
;

-- name: GetFeeds :many
select * from feeds
order by name asc
;

-- name: MarkFeedFetched :exec
update feeds
set updated_at = $2, 
    last_fetched_at = $2
where id = $1
;

-- name: GetNextFeedToFetch :one
select *
from feeds f
where f.id in (
    select feed_id
    from feed_follows ff
    where ff.user_id = $1
)
order by last_fetched_at asc nulls first
limit 1
;
