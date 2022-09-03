-- name: CreateAccount :one
INSERT INTO accounts(
  owner,
  balance,
  currency
) values (
  $1,$2,$3
) RETURNING *;


-- name: GetAccount :one
select * from accounts
where id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
select * from accounts
where id = $1 LIMIT 1
FOR NO KEY UPDATE;

  -- name: ListAccounts :many
  select * from accounts
  where owner = $1
  order by id
  LIMIT $2
  OFFSET $3;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
where id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
where id = sqlc.arg(id)
RETURNING *;


-- name: DeleteAccount :exec
DELETE FROM accounts 
where id =$1;

