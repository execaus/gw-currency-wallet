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
