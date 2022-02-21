-- +goose Up
-- +goose StatementBegin
CREATE TABLE "withdrawals" (
    "uuid" uuid not null primary key,
    "user_uuid" uuid not null,
    "order" varchar(100),
    "sum" numeric(13, 4) not null,
    "processed_at" timestamp null
);
CREATE INDEX ON "withdrawals" ("user_uuid");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "withdrawals";
-- +goose StatementEnd
