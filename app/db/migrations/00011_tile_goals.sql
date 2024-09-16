-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS public.tile_goals
(
    id serial,
    tile_id serial NOT NULL,
    goal_title character varying NOT NULL,
    target_value integer NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.tile_goals
    ADD FOREIGN KEY (tile_id)
    REFERENCES public.tiles (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE cascade
    NOT VALID;
END;

CREATE TABLE IF NOT EXISTS public.goal_progress
(
    id serial,
    goal_id serial NOT NULL,
    login_id serial NOT NULL,
    current_goal_progress integer NOT NULL,
    PRIMARY KEY (id)
);

ALTER TABLE IF EXISTS public.goal_progress
    ADD FOREIGN KEY (goal_id)
    REFERENCES public.tile_goals (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE cascade
    NOT VALID;
END;

ALTER TABLE IF EXISTS public.goal_progress
    ADD FOREIGN KEY (login_id)
    REFERENCES public.logins (id) MATCH SIMPLE
    ON UPDATE NO ACTION
    ON DELETE cascade
    NOT VALID;
END;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS public.goal_progress;
DROP TABLE IF EXISTS public.tile_goals;
-- +goose StatementEnd
