-- +goose Up
-- +goose StatementBegin
CREATE TABLE posts (
    id UUID primary key default gen_random_uuid(),
    created_at timestamp not null,
    updated_at timestamp,
    title varchar(150) not null,
    description text,
    url text unique not null,
    published_at timestamp,
    feed_id UUID not null,
    foreign key (feed_id) references feeds (id) on delete cascade
)
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
-- +goose StatementEnd
