-- +goose Up
-- +goose StatementBegin
create table outbox (
    id varchar(36) primary key not null default gen_random_uuid(),
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    "status" outbox_status_type not null default 'new',
    locked_to timestamp not null default now(),
    entity_id text,
    "payload" text,
    metadata jsonb not null default '{}'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists outbox;
-- +goose StatementEnd