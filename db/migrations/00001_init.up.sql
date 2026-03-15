CREATE TABLE IF NOT EXISTS users (
    idx INT PRIMARY KEY,
    uni_id INT UNIQUE,
    username VARCHAR(255) UNIQUE NOT NULL,

    name_ar VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20) NOT NULL,

    password VARCHAR(255) NOT NULL,
    verified tinyint NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active'
);

CREATE TABLE IF NOT EXISTS users_temp (
    idx INT PRIMARY KEY,
    uni_id INT UNIQUE,
    username VARCHAR(255) UNIQUE,

    name_ar VARCHAR(255),
    name_en VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20) UNIQUE,

    password VARCHAR(255),
    verified tinyint NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active'
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id INT NOT NULL,
    role VARCHAR(50) NOT NULL,
    INDEX (user_id),
    FOREIGN KEY (user_id) REFERENCES users(idx)
);

CREATE TABLE IF NOT EXISTS event (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    max_participants INT NOT NULL,
    outcomes TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL
);

CREATE TABLE IF NOT EXISTS event_component (
    id INT PRIMARY KEY AUTO_INCREMENT,
    event_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    max_score DECIMAL(10, 2) NOT NULL,

    FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS event_participant (
    id INT PRIMARY KEY AUTO_INCREMENT,
    event_id INT NOT NULL,
    user_id INT NOT NULL,
    grade DECIMAL(5, 2),
    status VARCHAR(20) NOT NULL,
    joined_at DATE NOT NULL,
    completed TINYINT NOT NULL,

    UNIQUE (event_id, user_id),

    FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(idx) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS component_score (
    id INT PRIMARY KEY AUTO_INCREMENT,
    participant_id INT NOT NULL,
    component_id INT NOT NULL,
    score DECIMAL(10, 2),

    UNIQUE (participant_id, component_id),

    FOREIGN KEY (participant_id) REFERENCES event_participant(id) ON DELETE CASCADE,
    FOREIGN KEY (component_id) REFERENCES event_component(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS store (
    id INT PRIMARY KEY AUTO_INCREMENT,
    fid VARCHAR(255) UNIQUE NOT NULL,
    mime VARCHAR(255) NOT NULL,
    size INT
);

CREATE TABLE IF NOT EXISTS certificate (
    id INT PRIMARY KEY AUTO_INCREMENT,
    cert_hash VARCHAR(255) UNIQUE NOT NULL,
    user_id INT NOT NULL,
    event_id INT NOT NULL,
    grade DECIMAL(5, 2) NOT NULL,
    issue_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL,

    UNIQUE (user_id, event_id),

    FOREIGN KEY (user_id) REFERENCES users(idx),
    FOREIGN KEY (event_id) REFERENCES event(id)
);

CREATE TABLE IF NOT EXISTS certificate_file (
    id INT PRIMARY KEY AUTO_INCREMENT,
    certificate_id INT NOT NULL,
    store_id INT NOT NULL UNIQUE,
    lang VARCHAR(10) NOT NULL,

    UNIQUE (certificate_id, lang),

    FOREIGN KEY (certificate_id) REFERENCES certificate(id),
    FOREIGN KEY (store_id) REFERENCES store(id) ON DELETE CASCADE
);
