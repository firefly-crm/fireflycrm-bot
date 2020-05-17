-- +goose Up
-- +goose StatementBegin
CREATE TABLE payments
(
    id              BIGSERIAL PRIMARY KEY,
    order_id        BIGINT,
    amount          INT         NOT NULL DEFAULT 0,
    payment_method  SMALLINT    NOT NULL,
    payment_link    TEXT        NOT NULL DEFAULT '',
    payed           BOOLEAN     NOT NULL DEFAULT FALSE,
    refunded        BOOLEAN     NOT NULL DEFAULT FALSE,
    refund_amount   INT         NOT NULL DEFAULT 0,
    bank_payment_id TEXT        NOT NULL DEFAULT '',
    expired         BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    payed_at        TIMESTAMPTZ,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_payed_amount() RETURNS TRIGGER AS
$update_payed_amount$
BEGIN
    IF (TG_OP = 'DELETE') THEN
        UPDATE
            orders
        SET payed_amount=COALESCE((SELECT SUM(amount) FROM payments WHERE order_id = OLD.order_id AND payed), 0),
            refund_amount=COALESCE((SELECT SUM(refund_amount) FROM payments WHERE order_id = OLD.order_id AND payed), 0)
        WHERE id = OLD.order_id;
        RETURN OLD;
    ELSE
        UPDATE
            orders
        SET payed_amount=COALESCE((SELECT SUM(amount) FROM payments WHERE order_id = NEW.order_id AND payed), 0),
            refund_amount=COALESCE((SELECT SUM(refund_amount) FROM payments WHERE order_id = OLD.order_id AND payed), 0)
        WHERE id = NEW.order_id;
        RETURN NEW;
    END IF;
    RETURN NULL;
END;
$update_payed_amount$ LANGUAGE plpgsql;

CREATE TRIGGER update_order_payed_amount
    AFTER INSERT OR UPDATE OR DELETE
    ON payments
    FOR EACH ROW
EXECUTE PROCEDURE update_payed_amount();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE payments;
DROP FUNCTION update_payed_amount();
DROP TRIGGER update_order_payed_amount ON payments;
-- +goose StatementEnd
