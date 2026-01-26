-- name: CreatePost :exec
insert into posts (
    id, 
    created_at, 
    updated_at, 
    title, 
    url, 
    description,
    published_at,
    feed_id
    )
values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
;

-- name: GetPostsForUser :many
select 
    p.title as post_title,
    p.description as post_description,
    p.published_at,
    f.name as feed_title
from posts p
join feed_follows ff on ff.feed_id = p.feed_id
join feeds f on f.id = ff.feed_id
where ff.user_id = $1
order by published_at desc
limit $2
;