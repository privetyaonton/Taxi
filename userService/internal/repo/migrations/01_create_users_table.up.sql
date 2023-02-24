CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    phone_number VARCHAR(30) NOT NULL,
    email VARCHAR(30) NOT NULL,
    password BYTEA NOT NULL,
    raiting FLOAT(8) NOT NULL,
    status INTEGER NOT NULL
);