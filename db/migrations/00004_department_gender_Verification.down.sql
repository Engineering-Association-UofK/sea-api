DROP EVENT IF EXISTS delete_expired_verification_codes;

ALTER TABLE users DROP COLUMN gender;
ALTER TABLE users DROP COLUMN department;