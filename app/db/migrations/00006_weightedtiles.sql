-- +goose Up
-- +goose StatementBegin
alter table public.tiles add column weight integer NOT NULL default 1;
alter table public.template_tiles add column weight integer NOT NULL default 1;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table public.tiles drop column weight;
alter table public.template_tiles drop column weight;
-- +goose StatementEnd
