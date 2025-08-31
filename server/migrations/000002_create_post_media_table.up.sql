CREATE TABLE post_media (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  post_id uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  url text NOT NULL,
  media_type varchar(32) NOT NULL, -- image, video, audio, file
  width int,
  height int,
  size_bytes bigint,
  meta jsonb DEFAULT '{}'::jsonb,
  sort_index int DEFAULT 0,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_post_media_post ON post_media(post_id);