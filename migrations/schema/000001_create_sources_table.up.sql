CREATE TABLE sources
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    feed_url   VARCHAR(255) NOT NULL,
    priority   INT          NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);
