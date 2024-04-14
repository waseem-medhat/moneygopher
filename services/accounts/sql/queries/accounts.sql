-- name: CreateAccount :one
INSERT INTO accounts ( id )
VALUES ( ? )
RETURNING *;
