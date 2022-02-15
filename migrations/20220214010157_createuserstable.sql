-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    "uuid" uuid not null PRIMARY KEY,
    "login" varchar(255) not null,
    "password" varchar(255) not null
);
CREATE UNIQUE INDEX ON users ("login");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table users;
-- +goose StatementEnd
