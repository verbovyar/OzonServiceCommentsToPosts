CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content VARCHAR(200) NOT NULL,
    author VARCHAR(200) NOT NULL,
    comments_enabled BOOLEAN, 
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    post_id VARCHAR(200) NOT NULL,
    parent_id VARCHAR(200) NOT NULL,
    author VARCHAR(200) NOT NULL,
    content VARCHAR(200) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
