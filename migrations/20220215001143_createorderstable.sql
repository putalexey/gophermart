-- +goose Up
-- +goose StatementBegin
CREATE TABLE "orders"
(
    uuid        uuid not null primary key,
    user_uuid   uuid not null,
    number      varchar(100),
    status      varchar(20),
    accrual     int,
    uploaded_at timestamp
);
CREATE INDEX ON orders ("user_uuid");
CREATE INDEX ON orders ("number");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "orders";
-- +goose StatementEnd
