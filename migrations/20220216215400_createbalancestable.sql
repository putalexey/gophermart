-- +goose Up
-- +goose StatementBegin
CREATE TABLE "balances"
(
    user_uuid uuid           not null primary key,
    current   numeric(13, 4) not null,
    withdrawn numeric(13, 4) not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "balances";
-- +goose StatementEnd
