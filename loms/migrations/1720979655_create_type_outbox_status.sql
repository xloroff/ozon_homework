-- +goose Up
-- +goose StatementBegin
create type outbox_status_type as enum (
    'new',
    'sent',
    'failed',
    'locked'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop type if exists outbox_status_type;
-- +goose StatementEnd