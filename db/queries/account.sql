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

-- name: ListAccount :many
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

-- name: DeleteAccount :exec
delete from accounts
where
  id = $1 
;