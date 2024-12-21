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
  username = $1
  hashed_password = $2
  full_name = $3
  email = $4
where
  id = $5
  returning *
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
-- name: DeleteUser :exec
delete from users
where
  id = $1 
;