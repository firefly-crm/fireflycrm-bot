-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id                  BIGINT      NOT NULL PRIMARY KEY,
    is_merchant         BOOLEAN     NOT NULL DEFAULT FALSE,
    active_order_msg_id BIGINT      NOT NULL DEFAULT 0,
    merchant_id         TEXT        NOT NULL DEFAULT '',
    secret_key          TEXT        NOT NULL DEFAULT '',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE customers
(
    id         BIGSERIAL   NOT NULL PRIMARY KEY,
    name       TEXT,
    email      TEXT,
    phone      TEXT,
    instagram  TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX customers_insta_idx ON customers (instagram);
CREATE UNIQUE INDEX customers_phone_idx ON customers (phone);
CREATE UNIQUE INDEX customers_email_idx ON customers (email);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE customers;
-- +goose StatementEnd
