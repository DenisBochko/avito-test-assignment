-- 000004_add_pull_requests_table.up.sql

CREATE TYPE pr_status AS ENUM ('OPEN','MERGED');

CREATE TABLE IF NOT EXISTS pull_requests (
    pull_request_id TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES users(id),
    status pr_status NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    merged_at TIMESTAMP WITH TIME ZONE NULL
);
