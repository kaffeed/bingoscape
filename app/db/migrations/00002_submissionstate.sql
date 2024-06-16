-- +goose Up
-- +goose StatementBegin
alter table public.submissions add column state integer NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table public.submissions drop column state;
-- +goose StatementEnd
