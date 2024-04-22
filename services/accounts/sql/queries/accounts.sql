-- name: CreateAccount :one
INSERT INTO accounts ( id, phone_number )
VALUES ( ?, ? )
RETURNING *;

-- name: GetAccountByID :one
SELECT *
FROM accounts
WHERE id = ?;
