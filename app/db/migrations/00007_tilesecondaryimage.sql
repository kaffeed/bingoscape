-- +goose Up
-- +goose StatementBegin
alter table public.tiles add column secondary_image_path character varying not null default 'https://i.ibb.co/7N9Pjcs/image.png';
alter table public.template_tiles add column secondary_image_path character varying not null default 'https://i.ibb.co/7N9Pjcs/image.png';;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table public.tiles drop column secondary_image_path;
alter table public.template_tiles drop column secondary_image_path;
-- +goose StatementEnd
