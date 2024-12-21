-- name: CreateUser :one
insert into
  users (owner, balance, currency)
values ($1, $2, $3)
  returning *
;

-- name: GetUser :one
select
  *
from
  users
where
  id = $1 
limit
  1
;

-- name: GetUserForUpdate :one
select
  *
from
  users
where
  id = $1 
limit
  1
for no key update
;

-- name: ListUsers :many
select
  *
from
  users
order by
  id
limit
  $1
offset
  $2
;

-- If there are no return values we use :exec instead of :one/many
-- name: UpdateUser :one
update users
set
  balance = $1 -- , another_param = $3
where
  id = $2
  returning *
;

-- name: UpdateUserBalance :one
update users
set
  balance = balance + sqlc.arg(amount) -- , another_param = $3
where
  id = sqlc.arg(id)
  returning *
;

-- name: DeleteUser :exec
delete from users
where
  id = $1 
;