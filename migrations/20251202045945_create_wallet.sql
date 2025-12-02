-- +goose Up
-- +goose StatementBegin

CREATE TABLE app.wallet (
    email VARCHAR(255) NOT NULL,
    currency VARCHAR(16) NOT NULL,
    balance FLOAT4 NOT NULL DEFAULT 0,
    PRIMARY KEY (email, currency)
);
ALTER TABLE app.wallet
    ADD CONSTRAINT account_balance_email_fk
    FOREIGN KEY (email) REFERENCES app.account(email) ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS app.account_balance;
-- +goose StatementEnd
