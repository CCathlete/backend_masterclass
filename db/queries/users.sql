-- name: CreateUser :one
insert into
  users (
    username, 
    hashed_password, 
    full_name,
    email
    ) values (
      $1, $2, $3, $4
    ) returning *
;

-- name: GetUser :one
select
  *
from
  users
where
  username = $1 
limit
  1
;

-- name: UpdateUser :one
update users
set
  username = sqlc.arg(new_username),
  hashed_password = sqlc.arg(hashed_password),
  full_name = sqlc.arg(full_name),
  email = sqlc.arg(email)
where
  username = sqlc.arg(username) -- The old username.
  returning *
;

-- name: ListUsers :many
select
  *
from
  users
order by
  username
limit
  $1
offset
  $2
;

-- If there are no return values we use :exec instead of :one/many
-- name: DeleteUser :exec
delete from users
where
  username = $1 
;