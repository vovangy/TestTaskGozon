CREATE TABLE IF NOT EXISTS comment_closure (
    ancestor_id BIGINT NULL REFERENCES comment(id)
    descendant_id BIGINT NOT NULL REFERENCES comment(id)
    PRIMARY KEY (ancestor_id, descendant_id)
);