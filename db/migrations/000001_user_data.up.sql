CREATE TABLE IF NOT EXISTS user_data (
    id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username TEXT CONSTRAINT username_length CHECK ( char_length(username) <= 40) NOT NULL,
    password_hash TEXT CONSTRAINT password_hash_length CHECK ( char_length(password_hash) <= 40) NOT NULL
);