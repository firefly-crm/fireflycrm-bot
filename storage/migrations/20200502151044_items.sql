-- +goose Up
-- +goose StatementBegin
CREATE TABLE items
(
    id      SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users,
    name    TEXT,
    type    SMALLINT
);
CREATE UNIQUE INDEX name_idx ON items (user_id, name);

CREATE TABLE receipt_items
(
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT                  NOT NULL DEFAULT '',
    item_id     INT REFERENCES items,
    order_id    INT REFERENCES orders NOT NULL,
    quantity    INT                   NOT NULL DEFAULT 1,
    price       INT                   NOT NULL DEFAULT 0,
    initialised BOOLEAN               NOT NULL DEFAULT FALSE
);
CREATE INDEX ord_idx ON receipt_items (order_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE items;
DROP TABLE receipt_items;
-- +goose StatementEnd
