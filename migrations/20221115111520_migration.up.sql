CREATE TABLE IF NOT EXISTS users
(
    id       serial PRIMARY KEY,
    login    VARCHAR(32)  not null,
    email    VARCHAR(254) not null,
    phone    bigint,
    password bytea
);

CREATE UNIQUE INDEX login_index ON users (login);
CREATE UNIQUE INDEX email_index ON users (email);