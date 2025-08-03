-- +goose Up
-- +goose StatementBegin
ALTER TABLE employee
    ADD is_deleted bool NOT NULL DEFAULT false;


ALTER TABLE role
    ADD is_deleted bool NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE employee
    DROP COLUMN is_deleted;

ALTER TABLE role
    DROP COLUMN is_deleted;
-- +goose StatementEnd
