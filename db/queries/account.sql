-- name: CreateAccount :one
INSERT INTO accounts(owner, balance, currency) values ($1, $2, $3) returning *;

-- name: GetAccount :one
SELECT * FROM accounts where id = $1 limit 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts where id = $1 limit 1 for no key update;

-- name: ListAccounts :many
select * from accounts order by id limit $1 offset $2;

-- name: UpdateAccount :one
update accounts set balance = $2 where id = $1 returning *;

-- name: AddAccountBalance :one
update accounts set balance = balance + sqlc.arg(amount) where id = sqlc.arg(id) returning *;

-- name: DeleteAccount :exec
delete from accounts where id = $1;