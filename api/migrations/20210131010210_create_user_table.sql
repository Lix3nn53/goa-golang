-- +goose Up
-- +goose StatementBegin

CREATE TABLE "users" (
     "id"         BIGSERIAL PRIMARY KEY,
     "name"       TEXT NOT NULL,
     "cif"       TEXT NOT NULL,
     "country"   TEXT NOT NULL,
     "postal_code" TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd