-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION tsm_system_rows;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP EXTENSION tsm_system_rows;
-- +goose StatementEnd
