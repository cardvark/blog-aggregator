-- +goose Up
-- +goose StatementBegin
CREATE TABLE feeds (
    id UUID primary key default gen_random_uuid(),
    created_at timestamp not null,
    updated_at timestamp,
    name varchar(150) not null,
    url text unique,
    user_id UUID    not null,
    foreign key (user_id) references users (id) on delete cascade
)
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE feeds;
-- +goose StatementEnd
