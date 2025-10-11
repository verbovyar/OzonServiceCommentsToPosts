CREATE TABLE IF NOT EXISTS posts (
    id VARCHAR(200) PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content VARCHAR(2000) NOT NULL,
    author VARCHAR(200) NOT NULL,
    comments_enabled BOOLEAN, 
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comments (
    id VARCHAR(200) PRIMARY KEY,
    post_id VARCHAR(200) NOT NULL,
    parent_id VARCHAR(200) NOT NULL,
    author VARCHAR(200) NOT NULL,
    content VARCHAR(2000) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_comments_parent_id ON comments(parent_id);