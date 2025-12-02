-- name: CreateAccount :one
INSERT INTO app.account (email, username, password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: IsAccountExistsByUsername :one
SELECT EXISTS (
    SELECT 1 FROM app.account WHERE username = $1
);

-- name: IsAccountExistsByEmail :one
SELECT EXISTS (
    SELECT 1 FROM app.account WHERE email = $1
);

-- name: GetAccountByUsername :one
SELECT *
FROM app.account
WHERE username = $1;

-- name: GetWalletsByEmail :many
SELECT *
FROM app.wallet
WHERE email = $1;

-- name: GetWalletForUpdate :one
SELECT *
FROM app.wallet
WHERE email = $1 and currency = $2
FOR UPDATE;

-- name: UpdateWallet :one
UPDATE app.wallet
SET balance = $3
WHERE email = $1 and currency = $2
RETURNING *;

-- name: IsExistCurrency :one
SELECT EXISTS (
    SELECT 1 FROM app.wallet WHERE email = $1 and currency = $2
);

-- name: CreateWallet :exec
INSERT INTO app.wallet (email, currency, balance)
VALUES ($1, $2, 0);