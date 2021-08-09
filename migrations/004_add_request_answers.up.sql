ALTER TABLE requests DROP COLUMN answer;
ALTER TABLE requests ADD COLUMN answer_text TEXT;
ALTER TABLE requests ADD COLUMN answered_by TEXT;
ALTER TABLE requests ADD COLUMN answered_at TIMESTAMP;
