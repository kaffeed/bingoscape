INSERT INTO logins (password, name, is_management)
VALUES (
    '$2a$08$voKjGNDQhECYiTpaJqx7CuVSeoVXNGSAArEb3PnfK1azcJGgR68EK',
    'test',
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

INSERT INTO bingos_logins (bingo_id, login_id)
VALUES (3,2);

--drop table schema_migrations;id:integer, login_id:integer