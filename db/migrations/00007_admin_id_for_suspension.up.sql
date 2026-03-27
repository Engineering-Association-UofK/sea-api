ALTER TABLE suspensions ADD COLUMN admin_id INT NOT NULL;
ALTER TABLE suspensions ADD CONSTRAINT fk_suspensions_users FOREIGN KEY (admin_id) REFERENCES users(idx);
ALTER TABLE suspension_history ADD COLUMN admin_id INT NOT NULL;
ALTER TABLE suspension_history ADD CONSTRAINT fk_suspension_history_users FOREIGN KEY (admin_id) REFERENCES users(idx);