-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id      BIGINT NOT NULL PRIMARY KEY,
    chat_id BIGINT NOT NULL
);
CREATE INDEX users_chat_id_idx ON users (chat_id);

CREATE TABLE customers
(
    id     BIGSERIAL NOT NULL PRIMARY KEY,
    name   TEXT,
    email  TEXT,
    phone  TEXT,
    social TEXT
);
CREATE INDEX customers_insta_idx ON customers (social);
CREATE INDEX customers_phone_idx ON customers (phone);
CREATE INDEX customers_email_idx ON customers (email);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE customers;
-- +goose StatementEnd
