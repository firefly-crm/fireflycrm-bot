-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_messages (
    id BIGINT NOT NULL,
    order_id BIGINT REFERENCES orders NOT NULL,
    user_id BIGINT REFERENCES users,
    display_mode SMALLINT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_messages;
-- +goose StatementEnd
