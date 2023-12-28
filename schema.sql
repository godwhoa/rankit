-- Users who participate in voting
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    display_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- Contests created by users
CREATE TABLE contests (
    id TEXT PRIMARY KEY,
    creator_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    content_type TEXT NOT NULL, -- v1/markdown, v1/image, v1/video, v1/audio, v1/link
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (creator_id) REFERENCES users(id)
);

-- Items that are being ranked
CREATE TABLE items (
    id TEXT PRIMARY KEY,
    contest_id TEXT NOT NULL,
    content_type TEXT NOT NULL, -- v1/markdown, v1/image, v1/video, v1/audio, v1/link
    content JSONB NOT NULL, -- { s3: "s3://" }
    elo_rating INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (contest_id) REFERENCES contests(id)
);

-- Votes table to store user votes
CREATE TABLE votes (
    id TEXT PRIMARY KEY,
    contest_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    winner_item_id TEXT NOT NULL,
    loser_item_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (contest_id) REFERENCES contests(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (winner_item_id) REFERENCES items(id),
    FOREIGN KEY (loser_item_id) REFERENCES items(id),
    UNIQUE (contest_id, user_id, winner_item_id, loser_item_id)
);

-- ELO rating history
CREATE TABLE elo_history (
    id TEXT PRIMARY KEY,
    item_id TEXT NOT NULL,
    elo_rating INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (item_id) REFERENCES items(id)
);