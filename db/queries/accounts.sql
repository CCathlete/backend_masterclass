-- name: CreateAccount :one
insert into
  accounts (owner, balance, currency)
values ($1, $2, $3)
  returning *
;

-- name: GetAccount :one
select
  *
from
  accounts
where
  id = $1 
limit
  1
;

-- name: GetAccountForUpdate :one
select
  *
from
  accounts
where
  id = $1 
limit
  1
for no key update
;

-- name: ListAccounts :many
select
  *
from
  accounts
order by
  id
limit
  $1
offset
  $2
;

-- If there are no return values we use :exec instead of :one/many
-- name: UpdateAccount :one
update accounts
set
  balance = $1 -- , another_param = $3
where
  id = $2
  returning *
;

-- name: UpdateAccountBalance :one
update accounts
set
  balance = balance + sqlc.arg(amount) -- , another_param = $3
where
  id = sqlc.arg(id)
  returning *
;

-- name: DeleteAccount :exec
delete from accounts
where
  id = $1 
;