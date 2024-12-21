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
  username = $2,
  hashed_password = $3,
  full_name = $4,
  email = $5
where
  username = $1 -- The old username.
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