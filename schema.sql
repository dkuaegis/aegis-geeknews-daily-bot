-- GeekNews RSS 데이터 저장을 위한 테이블 스키마
CREATE TABLE IF NOT EXISTS news (
    id SERIAL PRIMARY KEY,
    url VARCHAR(500) NOT NULL UNIQUE,
    title TEXT NOT NULL,
    author VARCHAR(255),
    content TEXT,
    published_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    sent BOOLEAN DEFAULT FALSE
);
