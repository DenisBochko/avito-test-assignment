-- 000001_add_teams_table.up.sql

CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    team_name TEXT NOT NULL UNIQUE
);
