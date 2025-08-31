CREATE TABLE follows (
  follower_id text NOT NULL,
  followee_id text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  status smallint NOT NULL DEFAULT 1, -- 1=active, 0=pending/blocked (for private accounts)
  PRIMARY KEY (follower_id, followee_id)
);

CREATE INDEX idx_follows_followee ON follows(followee_id);
CREATE INDEX idx_follows_follower ON follows(follower_id);
