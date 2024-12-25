-- name: CreateTransfer :one
insert into
  transfers (from_account_id, to_account_id, amount, currency)
values
  ($1, $2, $3, $4)
  returning *
;

-- name: GetTransfer :one
select
  *
from
  transfers
where
  id = $1 limit 1
  ;

-- name: GetTransfersFrom :many
select
  *
from
  transfers
where
  from_account_id = $1
;

-- name: GetTransfersTo :many
select
  *
from
  transfers
where
  to_account_id = $1
;

-- name: ListTransfers :many
select 
  t.*
from 
  transfers t
join 
  accounts a on t.from_account_id = a.id
where 
  a.owner = $1
order by 
  t.id
limit 
  $2
offset 
  $3
;

-- If there are no return values we use :exec instead of :one/many
-- name: UpdateTransfer :one
update transfers
set
  amount = $1 -- , another_param = $3
where
  id = $2
  returning *
;

-- name: DeleteTransfer :exec
delete from transfers
where
  id = $1;