-- +goose Up
-- +goose StatementBegin
alter table public.bingos add column codephrase character varying(255) NOT NULL DEFAULT 'Bingo!';
alter table public.bingos add column active boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table public.bingos drop column codephrase;
alter table public.bingos drop column active;
-- +goose StatementEnd
