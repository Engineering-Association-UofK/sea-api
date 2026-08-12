DROP TABLE logs;
DROP TABLE notifications;
DROP TABLE bot_command_options;
DROP TABLE bot_command_triggers;
DROP TABLE bot_command_translations;
DROP TABLE bot_commands;
DROP TABLE event_applications;
DROP TABLE event_form;

ALTER TABLE event DROP FOREIGN KEY event_ibfk_1;
ALTER TABLE event DROP COLUMN header_image_id;