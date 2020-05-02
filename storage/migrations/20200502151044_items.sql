-- +goose Up
-- +goose StatementBegin
CREATE TABLE items
(
    id   SERIAL PRIMARY KEY,
    name TEXT
);

CREATE TABLE receipt_items
(
    id SERIAL PRIMARY KEY,
    item_id INT REFERENCES items,
    bill_id INT REFERENCES bills
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE items;
DROP TABLE receipt_items;
-- +goose StatementEnd
