CREATE TABLE posts (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  author_id text NOT NULL,
  content text NOT NULL,                     -- raw markdown/plain text
  content_html text,                         -- pre-rendered sanitized HTML (optional)
  visibility smallint NOT NULL DEFAULT 0,    -- 0=public,1=followers,2=connections,3=private,4=direct
  reply_to_post_id uuid REFERENCES posts(id) ON DELETE SET NULL,
  original_post_id uuid REFERENCES posts(id) ON DELETE SET NULL, -- for reposts/shares
  language varchar(8),
  meta jsonb DEFAULT '{}'::jsonb,            -- link preview, client info, etc.
  likes_count bigint NOT NULL DEFAULT 0,
  comments_count bigint NOT NULL DEFAULT 0,
  shares_count bigint NOT NULL DEFAULT 0,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

-- tsvector for full-text search
ALTER TABLE posts ADD COLUMN IF NOT EXISTS search_vector tsvector;
CREATE INDEX idx_posts_search_vector ON posts USING GIN(search_vector);

-- trigger to update search_vector (simple example)
CREATE FUNCTION posts_search_vector_trigger() RETURNS trigger AS $$
begin
  new.search_vector :=
    to_tsvector('english', coalesce(new.content,'') || ' ' || coalesce(new.content_html,''));
  return new;
end
$$ LANGUAGE plpgsql;

CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE
  ON posts FOR EACH ROW EXECUTE PROCEDURE posts_search_vector_trigger();
