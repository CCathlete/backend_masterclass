-- name: CreateEntry :one
insert into entries (
  account_id,
  amount
)
values ($1, $2)
returning *
;

-- name: GetEntry :one
select * from entries
where id = $1 limit 1;

-- name: GetAccountEntries :many
select * from entries
where account_id = $1
order BY id
;

-- name: ListEntries :many
select * from entries
order BY id
limit $1
offset $2
;

-- If there are no return values we use :exec instead of :one/many
-- name: UpdateEntry :one
update entries set amount = $1 -- , another_param = $3
where id = $2
returning *
;

-- name: UpdateEntryByAccount :one
update entries set amount = $1 -- , another_param = $3
where account_id = $2
returning *
;

-- name: DeleteEntry :exec
delete from entries where id = $1;