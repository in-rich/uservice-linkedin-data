CREATE TABLE users (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    public_identifier VARCHAR(255) NOT NULL UNIQUE,

    first_name        VARCHAR(255) NOT NULL,
    last_name         VARCHAR(255) NOT NULL
);

--bun:split

CREATE INDEX users_public_identifier ON users(public_identifier);
