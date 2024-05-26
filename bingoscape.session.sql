INSERT INTO logins (password, name, is_management)
VALUES (
    '$2a$08$voKjGNDQhECYiTpaJqx7CuVSeoVXNGSAArEb3PnfK1azcJGgR68EK',
    'test',
    true
  );

INSERT INTO bingos (
    title,
    validFrom,
    validTo,
    size
  )
VALUES (
    'testbingo2',
    '2024-06-01',
    '2024-06-30',
    4
  );

  INSERT INTO bingos_logins (bingo_id, login_id)
  VALUES (
      2,
      4
    );  

SELECT b.id, title, "from", "to", size FROM bingos b
		JOIN bingos_logins bl ON b.id = bl.bingos_id
		JOIN logins l ON bl.logins_id = l.id
		WHERE b.id = 2