

-- name: CreateFeedFollow :one
With inserted_feed_follow as (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    returning *
)
select 
    iff.*,
    u.name as user_name,
    f.name as feed_name
from inserted_feed_follow iff
join users u on u.id = iff.user_id
join feeds f on f.id = iff.feed_id
;

-- name: GetFeedFollowsForUser :many
select 
    f.name as feed_name,
    u.name as user_name,
    f.url as url,
    fc.name as feed_creator
from feed_follows ff
join users u on u.id = ff.user_id
join feeds f on f.id = ff.feed_id
join users fc on f.user_id = fc.id
where ff.user_id = $1
order by f.name asc
;
