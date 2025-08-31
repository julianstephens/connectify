CREATE TABLE connections (
  user_a text NOT NULL,
  user_b text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  status smallint NOT NULL DEFAULT 1, -- 1=connected,0=pending
  PRIMARY KEY (user_a, user_b)
);
