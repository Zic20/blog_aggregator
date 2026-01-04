-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name VARCHAR(25) NOT NULL,
    CONSTRAINT uq_users_name UNIQUE(name)
);

-- +goose Down
DROP TABLE users;
