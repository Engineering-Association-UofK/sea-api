ALTER TABLE event ADD COLUMN header_image_id INT;
ALTER TABLE event ADD FOREIGN KEY (header_image_id) REFERENCES gallery_assets(id);

CREATE TABLE event_form (
    id INT PRIMARY KEY AUTO_INCREMENT,
    form_id INT NOT NULL,
    event_id INT NOT NULL,

    FOREIGN KEY (form_id) REFERENCES forms(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE
);

CREATE TABLE event_applications (
    id INT PRIMARY KEY AUTO_INCREMENT,
    event_id INT NOT NULL,
    user_id INT NOT NULL,
    status VARCHAR(20) NOT NULL,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- ------ BOT SCHEMA

CREATE TABLE bot_commands (
    id INT PRIMARY KEY AUTO_INCREMENT,
    keyword VARCHAR(50) NOT NULL,
    description VARCHAR(255)
);

CREATE TABLE bot_command_translations (
    command_id INT NOT NULL,
    text TEXT NOT NULL,
    language VARCHAR(10) NOT NULL,
    FOREIGN KEY (command_id) REFERENCES bot_commands(id)
);

CREATE TABLE bot_command_triggers (
    command_id INT NOT NULL,
    trigger_text VARCHAR(50) NOT NULL,
    language VARCHAR(10) NOT NULL,
    FOREIGN KEY (command_id) REFERENCES bot_commands(id)
);

CREATE TABLE bot_command_options (
    command_id INT NOT NULL,
    next_keyword VARCHAR(50) NOT NULL,
    FOREIGN KEY (command_id) REFERENCES bot_commands(id)
);

-- ------ NOTIFICATIONS SCHEMA

CREATE TABLE notifications (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    data JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_read TINYINT NOT NULL DEFAULT 0,

    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);

-- ------ LOGS SCHEMA

CREATE TABLE logs (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    action ENUM('CREATE', 'UPDATE', 'DELETE') NOT NULL,
    table_name VARCHAR(50) NOT NULL,
    object_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);

