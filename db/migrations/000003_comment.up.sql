CREATE TABLE IF NOT EXISTS comment (
    id BIGINT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    content TEXT CONSTRAINT content_length CHECK ( char_length(content) <= 2000) NOT NULL,
    user_id BIGINT NOT NULL REFERENCES user_data(id),
    post_id BIGINT NOT NULL REFERENCES post(id)
);