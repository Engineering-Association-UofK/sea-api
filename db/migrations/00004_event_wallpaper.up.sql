ALTER TABLE event ADD COLUMN wallpaper_id INT;
ALTER TABLE event ADD FOREIGN KEY (wallpaper_id) REFERENCES gallery_assets(id);