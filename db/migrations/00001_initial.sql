-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.logins
(
    id serial NOT NULL,
    name character varying(64) NOT NULL,
    is_management boolean NOT NULL DEFAULT 'false',
    password character varying NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.bingos
(
    id serial NOT NULL,
    title character varying(255) NOT NULL,
    validFrom timestamptz NOT NULL,
    validTo timestamptz NOT NULL,
    rows integer,
    cols integer,
    description character varying(5000),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.bingos_logins
(
    bingo_id serial NOT NULL,
    login_id serial NOT NULL
);

CREATE TABLE IF NOT EXISTS public.tiles
(
    id serial NOT NULL,
    title character varying(250) NOT NULL,
    imagepath character varying NOT NULL,
    description character varying(5000) NOT NULL,
    bingo_id serial NOT NULL,
    PRIMARY KEY (id)
);
    
CREATE TABLE IF NOT EXISTS public.template_tiles
(
    id serial NOT NULL,
    title character varying(250) NOT NULL,
    imagepath character varying NOT NULL,
    description character varying(5000) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.submissions
(
    id serial NOT NULL,
    login_id serial NOT NULL,
    tile_id serial NOT NULL,
    date timestamptz NOT NULL,
    comment character varying(5000),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.submission_images
(
    id serial,
    path character varying NOT NULL,
    submission_id serial NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.bingos_logins
    ADD FOREIGN KEY (bingo_id)
    REFERENCES public.bingos (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE cascade
    NOT VALID;


ALTER TABLE IF EXISTS public.bingos_logins
    ADD FOREIGN KEY (login_id)
    REFERENCES public.logins (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public.tiles
    ADD FOREIGN KEY (bingo_id)
    REFERENCES public.bingos (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE CASCADE
    NOT VALID;


ALTER TABLE IF EXISTS public.submissions
    ADD FOREIGN KEY (login_id)
    REFERENCES public.logins (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public.submissions
    ADD FOREIGN KEY (tile_id)
    REFERENCES public.tiles (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;


ALTER TABLE IF EXISTS public.submission_images
    ADD FOREIGN KEY (submission_id)
    REFERENCES public.submissions (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table public.submission_images;
drop table public.submissions;
drop table public.tiles;
drop table public.template_tiles;
drop table public.bingos_logins;
drop table public.logins;
drop table public.bingos;
-- +goose StatementEnd
