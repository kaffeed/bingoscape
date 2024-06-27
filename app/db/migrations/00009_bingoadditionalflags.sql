-- +goose Up
-- +goose StatementBegin
alter table public.bingos add column submissions_closed bool not null default false;
alter table public.bingos add column leaderboard_public bool not null default false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table public.bingos drop column submissions_closed;
alter table public.bingos drop column leaderboard_public;
-- +goose StatementEnd
