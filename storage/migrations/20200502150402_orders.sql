-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders
(
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT REFERENCES users,
    customer_id     BIGINT REFERENCES customers,
    message_id      BIGINT,
    description     TEXT     NOT NULL DEFAULT '',
    amount          INT      NOT NULL DEFAULT 0,
    active_item_id  BIGINT REFERENCES receipt_items,
    hint_message_id BIGINT,
    state           SMALLINT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
