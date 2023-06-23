CREATE TYPE role AS ENUM ('user', 'editor', 'admin');
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(16) UNIQUE,
    email VARCHAR(127) UNIQUE,
    pass_hash VARCHAR(60) NOT NULL,
    role role NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP 
);
---- create above / drop below ----

DROP TABLE users;
DROP TYPE role;
