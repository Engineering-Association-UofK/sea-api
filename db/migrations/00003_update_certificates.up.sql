ALTER TABLE certificate ADD COLUMN type VARCHAR(255) DEFAULT 'participation';
ALTER TABLE certificate ADD COLUMN cert_version VARCHAR(255) DEFAULT 'v0.1';

ALTER TABLE collaborators ADD COLUMN title_ar VARCHAR(255) DEFAULT 'المقدم';
ALTER TABLE collaborators ADD COLUMN title_en VARCHAR(255) DEFAULT 'Presenter';

-- After adding default values to existing record, turn them to NOT NULL

ALTER TABLE certificate MODIFY COLUMN type VARCHAR(255) NOT NULL;
ALTER TABLE certificate MODIFY COLUMN cert_version VARCHAR(255) NOT NULL;

ALTER TABLE collaborators MODIFY COLUMN title_ar VARCHAR(255) NOT NULL;
ALTER TABLE collaborators MODIFY COLUMN title_en VARCHAR(255) NOT NULL;