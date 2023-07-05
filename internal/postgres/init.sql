CREATE USER docker WITH PASSWORD 'docker';

CREATE DATABASE certificate_dev;

GRANT ALL PRIVILEGES ON DATABASE certificate_dev TO docker;

\connect certificate_dev;

CREATE EXTENSION pgcrypto;

CREATE TABLE users (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    password TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    active BOOL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE certificates (
    uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_uuid UUID REFERENCES users(uuid),
    private_key VARCHAR NOT NULL,
    body VARCHAR NOT NULL,
    active BOOL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX user_idx ON certificates (user_uuid, active);

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO docker;