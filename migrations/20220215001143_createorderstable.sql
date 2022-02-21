-- +goose Up
-- +goose StatementBegin
CREATE TABLE "orders"
(
    uuid        uuid not null primary key,
    user_uuid   uuid not null,
    number      varchar(100) not null,
    status      varchar(20) not null,
    accrual     numeric(13, 4) not null,
    uploaded_at timestamp not null
);
CREATE INDEX ON orders ("user_uuid");
CREATE INDEX ON orders ("number");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "orders";
-- +goose StatementEnd
