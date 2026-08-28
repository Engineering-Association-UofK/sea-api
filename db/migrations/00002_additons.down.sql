DROP TABLE IF EXISTS feedback;
DROP TABLE IF EXISTS logs;
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS bot_user_states;
DROP TABLE IF EXISTS bot_edge_translations;
DROP TABLE IF EXISTS bot_actions;
DROP TABLE IF EXISTS bot_edges;
DROP TABLE IF EXISTS bot_node_translations;
DROP TABLE IF EXISTS bot_nodes;
DROP TABLE IF EXISTS event_applications;
DROP TABLE IF EXISTS event_form;

ALTER TABLE event DROP FOREIGN KEY fk_event_header_image;
ALTER TABLE event DROP COLUMN header_image_id;