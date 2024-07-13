-- +goose Up
CREATE TABLE feed_follows (
    id UUID PRIMARY KEY,
    feed_id UUID NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_uid FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_fid FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE,
    CONSTRAINT uk_user_feed UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;