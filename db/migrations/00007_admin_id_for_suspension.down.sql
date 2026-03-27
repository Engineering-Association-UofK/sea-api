ALTER TABLE suspensions DROP CONSTRAINT fk_suspensions_users;
ALTER TABLE suspensions DROP COLUMN admin_id;
ALTER TABLE suspension_history DROP CONSTRAINT fk_suspension_history_users;
ALTER TABLE suspension_history DROP COLUMN admin_id;