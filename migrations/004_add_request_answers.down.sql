ALTER TABLE requests DROP COLUMN answer_text;
ALTER TABLE requests DROP COLUMN answered_by;
ALTER TABLE requests DROP COLUMN answered_at;
ALTER TABLE requests ADD COLUMN answer TEXT NOT NULL;
