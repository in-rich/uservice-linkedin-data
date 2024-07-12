CREATE TABLE companies (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    public_identifier VARCHAR(255) NOT NULL UNIQUE,

    name        VARCHAR(255) NOT NULL
);

--bun:split

CREATE INDEX companies_public_identifier ON companies(public_identifier);
