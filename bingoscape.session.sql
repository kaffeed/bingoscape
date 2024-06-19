INSERT INTO logins (password, name, is_management)
VALUES (
    '$2a$08$voKjGNDQhECYiTpaJqx7CuVSeoVXNGSAArEb3PnfK1azcJGgR68EK',
    'Major',
    true
  );

INSERT INTO bingos (title, validfrom, validto, rows, cols)
VALUES (
    'Testbingo 2',
    '2024-06-01',
    '2024-06-30',
    5,
    4 
  );

  INSERT INTO bingos_logins (bingo_id, login_id)
  VALUES (3, 2);

SELECT b.id, b.title, b.validFrom, b.validTo, b.rows, b.cols FROM bingos b
		JOIN bingos_logins bl ON b.id = bl.bingo_id
		JOIN logins l ON bl.login_id = l.id
WHERE l.id = 1 

SELECT l.Id, l.name FROM public.logins l
	JOIN bingos_logins bl ON l.id = bl.login_id
	WHERE bl.bingo_id = 1 

SELECT l.id, l.name FROM public.logins l
	WHERE l.id NOT IN (SELECT login_id from public.bingos_logins WHERE bingo_id = 1) and not l.is_management;

INSERT INTO bingos_logins (bingo_id, login_id)
VALUES (1,1);

--drop table schema_migrations;id:integer, login_id:integer



SELECT s.id, s.login_id, s.tile_id, s.date, s.state
	FROM public.submissions s 
	WHERE s.id = 1

begin;
drop table public.submission_images;
drop table public.submissions;
drop table public.tiles;
drop table public.template_tiles;
drop table public.bingos_logins;
drop table public.logins;
drop table public.bingos;
drop table public.goose_db_version;
commit;

SELECT s.id, s.login_id, l.name, s.tile_id, s.date, s.comment 
	FROM public.Submissions s 
	JOIN public.logins l ON l.id = s.login_id
	WHERE tile_id = 1;

SELECT path FROM public.submission_images WHERE submission_id = 38;

SELECT t.id, t.imagepath, t.description, t.bingo_id, s.id as submission_id, s.tile_id, s.date, s.login_id
	FROM tiles t 
JOIN submissions s ON s.tile_id = t.id
	WHERE bingo_id = 1 ORDER BY t.id ASC;
