-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users
(
    id           UUID                 DEFAULT uuid_generate_v4() PRIMARY KEY,
    auth_user_id UUID        NOT NULL UNIQUE,
    login        VARCHAR(50) NOT NULL UNIQUE,
    created_at   TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE rooms
(
    id         UUID                  DEFAULT uuid_generate_v4() PRIMARY KEY,
    name       VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users_rooms
(
    id        UUID               DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id   UUID      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    room_id   UUID      NOT NULL REFERENCES rooms (id) ON DELETE CASCADE,
    joined_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, room_id)
);

CREATE TABLE messages
(
    id         UUID               DEFAULT uuid_generate_v4() PRIMARY KEY,
    room_id    UUID      NOT NULL REFERENCES rooms (id) ON DELETE CASCADE,
    user_id    UUID      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    content    TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE users_rooms;
DROP TABLE rooms;
DROP TABLE users;
DROP TABLE messages;