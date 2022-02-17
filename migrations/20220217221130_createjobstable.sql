-- +goose Up
-- +goose StatementBegin
CREATE TABLE "jobs" (
    "uuid" uuid not null primary key,
    "order_uuid" uuid not null,
    "proceed_at" timestamp not null,
    "tries" int not null default 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "jobs";
-- +goose StatementEnd
