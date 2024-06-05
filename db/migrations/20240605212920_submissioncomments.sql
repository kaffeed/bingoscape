-- +goose Up
-- +goose StatementBegin
alter table public.submissions drop column comment;

CREATE TABLE IF NOT EXISTS public.submission_comments
(
    id serial,
    comment character varying NOT NULL,
    submission_id serial NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now() at time zone 'utc'),
    login_id serial NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.submission_comments
    ADD FOREIGN KEY (submission_id)
    REFERENCES public.submissions (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;
END;

ALTER TABLE IF EXISTS public.submission_comments
    ADD FOREIGN KEY (login_id)
    REFERENCES public.logins (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE NO ACTION
    NOT VALID;
END;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table public.submissions add column comment character varying;
drop table public.submission_comments;
-- +goose StatementEnd
