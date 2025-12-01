-- +goose Up
-- +goose StatementBegin
CREATE TABLE app.account (
    email VARCHAR(255) PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS app.account;
-- +goose StatementEnd
