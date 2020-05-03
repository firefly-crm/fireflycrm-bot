-- +goose Up
-- +goose StatementBegin
CREATE TABLE bills
(
    id           BIGSERIAL PRIMARY KEY,
    order_id     BIGINT REFERENCES orders,
    amount       INT,
    payment_link TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bills;
-- +goose StatementEnd
