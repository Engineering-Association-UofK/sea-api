ALTER TABLE event ADD COLUMN header_image_id INT;
ALTER TABLE event ADD CONSTRAINT fk_event_header_image FOREIGN KEY (wallpaper_id) REFERENCES gallery_assets(id);

ALTER TABLE event DROP FOREIGN KEY fk_event_wallpaper;
ALTER TABLE event DROP COLUMN wallpaper_id;