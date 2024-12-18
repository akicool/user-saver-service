CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    name VARCHAR(40),
    email VARCHAR(255),
    password VARCHAR(255)
);