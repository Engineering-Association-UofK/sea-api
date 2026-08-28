ALTER TABLE event ADD COLUMN wallpaper_id INT;
ALTER TABLE event ADD CONSTRAINT fk_event_wallpaper FOREIGN KEY (wallpaper_id) REFERENCES gallery_assets(id);

ALTER TABLE event DROP FOREIGN KEY fk_event_header_image;
ALTER TABLE event DROP COLUMN header_image_id;