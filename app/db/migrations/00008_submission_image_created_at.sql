-- +goose Up
-- +goose StatementBegin
alter table public.submission_images add column created_at timestamp default now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table public.submission_images drop column created_at;
-- +goose StatementEnd
