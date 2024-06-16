-- +goose Up
-- +goose StatementBegin
CREATE TYPE SUBMISSIONSTATE AS ENUM (
  'Submitted',
  'ActionRequired',
  'Accepted'
);

alter table public.submissions drop column state;
alter table public.submissions add column state SUBMISSIONSTATE NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

alter table public.submissions drop column state;
alter table public.submissions add column state integer default 0;

drop type SUBMISSIONSTATE;
-- +goose StatementEnd
