-- name: CreateAccount :one
insert into accounts (
  owner,
  balance,
  currency
)
values ($1, $2, $3)
returning *
;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 limit 1;

-- name: ListAccount :many
SELECT * FROM accounts
ORDER BY id
limit $1
offset $2;