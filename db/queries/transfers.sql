-- name: CreateTransfer :one
insert into
  transfers (from_account_id, to_account_id, amount)
values
  ($1, $2, $3)
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
  *
from
  transfers
order by
  id
limit
  $1
offset
  $2
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