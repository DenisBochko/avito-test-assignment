-- 000003_add_team_lnk_table.up.sql

CREATE TABLE IF NOT EXISTS team_lnk (
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE
);
