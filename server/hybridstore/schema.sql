CREATE TABLE IF NOT EXISTS notes (
  id TEXT NOT NULL PRIMARY KEY,
  path TEXT NOT NULL UNIQUE,
  title TEXT NOT NULL,
  sourceURL TEXT NOT NULL,
  source TEXT NOT NULL,
  format TEXT NOT NULL,
  created INTEGER NOT NULL,
  updated INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS tags_notes (
  tag TEXT NOT NULL,
  note_id TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_tags_notes_u ON tags_notes (tag,note_id);

CREATE INDEX IF NOT EXISTS idx_tags_notes_t ON tags_notes (tag);

CREATE INDEX IF NOT EXISTS idx_tags_notes_a ON tags_notes (note_id);
