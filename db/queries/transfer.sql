-- name: CreateTransfer :one
insert into transfers (from_account_id, to_account_id, amount) values ($1, $2, $3) returning *;

-- name: GetTransfer :one
select * from transfers where id = $1 limit 1;

-- name: ListTransfersToAccount :many
select * from transfers where to_account_id = $1 order by id limit $2 offset $3;

-- name: ListTransfersFromAccount :many
select * from transfers where from_account_id = $1 order by id limit $2 offset $3;
