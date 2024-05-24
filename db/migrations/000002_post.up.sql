CREATE TABLE IF NOT EXISTS post (
    id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES user_data(id),
    is_commented BOOLEAN NOT NULL DEFAULT true,
    title TEXT CONSTRAINT title_length CHECK ( char_length(title) <= 150) NOT NULL,
    content TEXT CONSTRAINT content_length CHECK ( char_length(content) <= 5000) NOT NULL
);