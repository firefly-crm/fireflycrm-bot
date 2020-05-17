-- +goose Up
-- +goose StatementBegin
CREATE TABLE items
(
    id         SERIAL PRIMARY KEY,
    user_id    BIGINT REFERENCES users,
    name       TEXT,
    type       SMALLINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX name_idx ON items (user_id, name);

CREATE TABLE receipt_items
(
    id           BIGSERIAL PRIMARY KEY,
    name         TEXT                  NOT NULL DEFAULT '',
    item_id      INT REFERENCES items,
    order_id     INT REFERENCES orders NOT NULL,
    quantity     INT                   NOT NULL DEFAULT 1,
    price        INT                   NOT NULL DEFAULT 0,
    payed_amount INT                   NOT NULL DEFAULT 0,
    initialised  BOOLEAN               NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMPTZ           NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE orders
    ADD CONSTRAINT receipt_items_fk FOREIGN KEY (active_item_id) REFERENCES receipt_items (id);
CREATE INDEX ord_idx ON receipt_items (order_id);

CREATE OR REPLACE FUNCTION update_order_amount() RETURNS TRIGGER AS
$update_order_amount$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        UPDATE
            orders
        SET amount=COALESCE((SELECT SUM(price * quantity) FROM receipt_items WHERE order_id = OLD.order_id), 0)
        WHERE id = OLD.order_id;

        RETURN OLD;
    ELSE
        UPDATE
            orders
        SET amount=COALESCE((SELECT SUM(price * quantity) FROM receipt_items WHERE order_id = NEW.order_id), 0)
        WHERE id = NEW.order_id;

        RETURN NEW;
    END IF;

    RETURN NULL;
END;
$update_order_amount$ LANGUAGE plpgsql;

CREATE TRIGGER update_order_amount
    AFTER INSERT OR UPDATE OR DELETE
    ON receipt_items
    FOR EACH ROW
EXECUTE PROCEDURE update_order_amount();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE items;
DROP TABLE receipt_items;
DROP FUNCTION IF EXISTS update_order_amount();
DROP TRIGGER IF EXISTS update_order_amount ON receipt_items;
-- +goose StatementEnd
