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

CREATE TABLE bot_nodes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    node_type VARCHAR(20) DEFAULT 'message', -- 'message', 'input', 'action'
    is_start BOOLEAN DEFAULT FALSE,          -- True for the Welcome node only
    is_locked BOOLEAN DEFAULT FALSE, 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE bot_node_translations (
    node_id INT NOT NULL,
    language VARCHAR(5) NOT NULL,            -- 'en', 'ar'
    content TEXT NOT NULL,
    PRIMARY KEY (node_id, language),
    FOREIGN KEY (node_id) REFERENCES bot_nodes(id) ON DELETE CASCADE
);

CREATE TABLE bot_edges (
    id INT AUTO_INCREMENT PRIMARY KEY,
    from_node_id INT NOT NULL,
    to_node_id INT NOT NULL,
    keyword VARCHAR(50) NOT NULL,           
    FOREIGN KEY (from_node_id) REFERENCES bot_nodes(id) ON DELETE CASCADE,
    FOREIGN KEY (to_node_id) REFERENCES bot_nodes(id) ON DELETE CASCADE
);

CREATE TABLE bot_actions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    node_id INT NOT NULL UNIQUE,
    action_type VARCHAR(100) NOT NULL,
    action_text TEXT NOT NULL,
    FOREIGN KEY (node_id) REFERENCES bot_nodes(id) ON DELETE CASCADE
);

CREATE TABLE bot_edge_translations (
    edge_id INT NOT NULL,
    language VARCHAR(5) NOT NULL,
    label VARCHAR(255) NOT NULL,
    PRIMARY KEY (edge_id, language),
    FOREIGN KEY (edge_id) REFERENCES bot_edges(id) ON DELETE CASCADE
);

CREATE TABLE bot_user_states (
    session_id VARCHAR(255) PRIMARY KEY,
    current_node_id INT NOT NULL,
    user_id BIGINT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (current_node_id) REFERENCES bot_nodes(id) ON DELETE CASCADE
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

CREATE TABLE feedback (
    id INT PRIMARY KEY AUTO_INCREMENT,
    message TEXT NOT NULL,
    user_id INT,
    type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);
