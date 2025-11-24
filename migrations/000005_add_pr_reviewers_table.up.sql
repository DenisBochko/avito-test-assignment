-- 000005_add_pr_reviewers_table.up.sql

CREATE TABLE IF NOT EXISTS pr_reviewers (
    pull_request_id TEXT REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    reviewer_id TEXT REFERENCES users(id),
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    PRIMARY KEY(pull_request_id, reviewer_id)
);
