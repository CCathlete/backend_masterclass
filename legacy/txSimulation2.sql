-- Tx 1 => 10 from acc1 to acc2.
begin;

-- name: UpdateAccount :one
update accounts
set
  balance = balance - 10
where
  id = 1
  returning *
;
-- name: UpdateAccount :one
update accounts
set
  balance = balance + 10
where
  id = 2
  returning *
;

rollback;

-- Tx 2 => 10 from acc2 to acc1.
begin;

-- name: UpdateAccount :one
update accounts
set
  balance = balance - 10
where
  id = 2
  returning *
;
-- name: UpdateAccount :one
update accounts
set
  balance = balance + 10
where
  id = 1
  returning *
;

rollback;