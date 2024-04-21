-- name: CreateAccount :one
INSERT INTO accounts ( id )
VALUES ( ? )
RETURNING *;

-- name: GetAccountByID :one
SELECT *
FROM accounts
WHERE id = ?;
