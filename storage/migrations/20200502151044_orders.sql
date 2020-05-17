-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders
(
    id                BIGSERIAL PRIMARY KEY,
    user_id           BIGINT REFERENCES users,
    customer_id       BIGINT REFERENCES customers,
    description       TEXT        NOT NULL DEFAULT '',
    amount            INT         NOT NULL DEFAULT 0,
    payed_amount      INT         NOT NULL DEFAULT 0,
    refund_amount     INT         NOT NULL DEFAULT 0,
    active_item_id    BIGINT,
    active_payment_id BIGINT      REFERENCES payments ON DELETE SET NULL,
    hint_message_id   BIGINT,
    order_state       SMALLINT    NOT NULL DEFAULT 0,
    edit_state        SMALLINT    NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE payments
    ADD CONSTRAINT orders_fk FOREIGN KEY (order_id) REFERENCES orders (id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
