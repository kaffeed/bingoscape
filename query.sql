-- name: GetAllLogins :many
select * from public.logins;

-- name: CreateLogin :exec
INSERT INTO logins(password, name, is_management) VALUES($1, $2, $3) returning *;