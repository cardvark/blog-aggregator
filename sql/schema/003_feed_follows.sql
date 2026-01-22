-- +goose Up
-- +goose StatementBegin
CREATE TABLE feed_follows (
    id UUID primary key default gen_random_uuid(),
    created_at timestamp not null,
    updated_at timestamp,
    user_id UUID    not null,
    feed_id UUID not null,
    foreign key (user_id) references users (id) on delete cascade,
    foreign key (feed_id) references feeds (id) on delete cascade,
    unique (user_id, feed_id)
)
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feed_folows;
-- +goose StatementEnd
