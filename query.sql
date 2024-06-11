-- name: GetAllLogins :many
select * from public.logins;

-- name: CreateLogin :exec
INSERT INTO logins(password, name, is_management) VALUES($1, $2, $3) returning *;

-- name: GetLoginByName :one
SELECT * FROM logins WHERE name = $1;

-- name: DeleteLogin :exec
DELETE FROM logins WHERE id = $1;

-- name: UpdateLoginPassword :one
UPDATE logins SET password = $2 WHERE id = $1 returning *;

-- name: CreateSubmissionComment :exec
INSERT INTO public.submission_comments (submission_id, login_id, comment) values ($1, $2, $3);

-- name: GetCommentsForSubmission :many
SELECT c.submission_id, l.name as managementuser, c.comment, c.created_at FROM public.submission_comments c
JOIN public.logins l on c.login_id = l.id
WHERE c.submission_id = $1;

-- name: UpdateSubmissionState :one
UPDATE public.submissions SET state = $2 WHERE id = $1 returning *;

-- name: GetImagesForSubmission :many
SELECT path FROM public.submission_images WHERE submission_id = $1;

-- name: GetSubmissionsForTile :many
SELECT sqlc.embed(Submissions), sqlc.embed(logins)
FROM submissions
JOIN logins ON logins.id = submissions.login_id
WHERE submissions.tile_id = $1;

-- name: GetSubmissionsForTileAndLogin :many
SELECT sqlc.embed(submissions), sqlc.embed(logins)
FROM submissions
JOIN logins ON logins.id = submissions.login_id
WHERE submissions.tile_id = $1 AND submissions.login_id = $2;

-- name: DeleteBingoById :exec
DELETE FROM bingos WHERE id = $1;

-- name: GetSubmissionIdForTileAndLogin :one 
SELECT id FROM public.submissions WHERE tile_id = $1 AND login_id = $2;

-- name: CreateSubmission :one
INSERT INTO public.submissions (login_id, tile_id, state) values ($1, $2, $3) returning *;

-- name: CreateSubmissionImage :exec 
INSERT INTO submission_images(path, submission_id) VALUES ($1, $2); 

-- name: CreateTemplateTile :one
INSERT INTO template_tiles(title, imagepath, description) VALUES ($1, $2, $3) returning *;

-- name: GetTemplateTiles :many
SELECT * FROM template_tiles;

-- name: GetPossibleBingoParticipants :many
SELECT l.id, l.name FROM public.logins l
	WHERE l.id NOT IN (SELECT login_id from public.bingos_logins WHERE bingo_id = $1)
	AND not l.is_management;

-- name: GetBingoParticipants :many 
SELECT l.Id, l.name FROM public.logins l
	JOIN bingos_logins bl ON l.id = bl.login_id
	WHERE bl.bingo_id = $1;

-- name: GetBingos :many
SELECT * FROM bingos;

-- name: GetBingosForLogin :many
SELECT b.* FROM bingos b
JOIN bingos_logins bl ON b.id = bl.bingo_id
JOIN logins l ON bl.login_id = l.id
WHERE l.id = $1;

-- name: GetBingoById :one
SELECT * FROM bingos WHERE id = $1;

-- name: GetSubmissionById :one
SELECT * FROM submissions WHERE id = $1;

-- name: UpdateTile :one
UPDATE tiles SET title = $2, imagepath = $3, description = $4 WHERE id = $1 returning *;

-- name: GetTileById :one
SELECT * FROM tiles WHERE id = $1;

-- name: GetTilesForBingo :many
SELECT *
FROM tiles 
WHERE bingo_id = $1 ORDER BY id ASC;

-- name: DeleteBingoParticipant :exec
DELETE FROM bingos_logins WHERE login_id = $1 AND bingo_id = $2;

-- name: CreateBingoParticipant :exec
INSERT INTO bingos_logins (login_id, bingo_id) VALUES ($1, $2);

-- name: CreateBingo :one
INSERT INTO bingos (title, validFrom, validTo, rows, cols, description, ready, codephrase) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: CreateTile :one
INSERT INTO tiles(title, imagepath, description, bingo_id) VALUES ($1, $2, $3, $4) returning *;
