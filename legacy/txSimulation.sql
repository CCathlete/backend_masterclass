begin;

-- name: CreateTransfer :one
insert into
  transfers (from_account_id, to_account_id, amount)
values
  (1, 2, 10)
  returning *
;

-- name: CreateEntry :one
insert into entries (
  account_id,
  amount
)
values (1, -10)
returning *
;

-- name: CreateEntry :one
insert into entries (
  account_id,
  amount
)
values (2, 10)
returning *
;

-- name: GetAccountForUpdate :one
select
  *
from
  accounts
where
  id = 1 
limit
  1
for update
;

-- name: UpdateAccount :one
update accounts
set
  balance = 90 -- , another_param = $3
where
  id = 1
  returning *
;

-- name: GetAccountForUpdate :one
select
  *
from
  accounts
where
  id = 2 
limit
  1
for update
;

-- name: UpdateAccount :one
update accounts
set
  balance = 110 -- , another_param = $3
where
  id = 2
  returning *
;

commit;