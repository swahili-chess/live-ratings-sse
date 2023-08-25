-- name: GetallStatic :many
SELECT id,lichess_id,username,rapid FROM static;


-- name: GetallDynamic :many
SELECT id,lichess_id,username,rapid FROM dynamic;


-- name: InsertStatic :exec
INSERT INTO static (lichess_id,username,rapid) VALUES ($1,$2,$3);

-- name: InsertDynamic :exec
INSERT INTO dynamic (lichess_id,username,rapid) VALUES ($1,$2,$3);


-- name: UpdateDynamic :exec
UPDATE dynamic
SET rapid = $1, modified_at= NOW()
WHERE lichess_id = $2;

-- name: CheckEntriesStatic :one
SELECT COUNT(*) FROM static;

-- name: CheckEntriesDynamic :one
SELECT COUNT(*) FROM dynamic;