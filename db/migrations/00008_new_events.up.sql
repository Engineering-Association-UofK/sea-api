DROP TABLE event_applications;
DROP TABLE event_form;
DROP TABLE certificate_file;
DROP TABLE certificate;
DROP TABLE component_score;
DROP TABLE event_participant;
DROP TABLE event_component;
DROP TABLE event;
DROP TABLE collaborators;

CREATE TABLE events (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    background_id INT NOT NULL,
    
    belonging VARCHAR(255) NOT NULL,

    require_applying TINYINT NOT NULL,
    form_id INT NULL,
    max_applications INT NULL,
    
    created_at DATE NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,

    FOREIGN KEY (form_id) REFERENCES forms(id) ON DELETE CASCADE
);

CREATE TABLE event_coords (
    id INT PRIMARY KEY AUTO_INCREMENT,
    event_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL,

    FOREIGN KEY (event_id) REFERENCES events(id) ON UPDATE CASCADE
);

CREATE TABLE event_participation (
    event_id INT NOT NULL,
    user_id INT NOT NULL,
    joined_at DATE NOT NULL,

    UNIQUE KEY (user_id, event_id),

    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events(id) ON UPDATE CASCADE
);

CREATE TABLE event_application (
    event_id INT NOT NULL,
    user_id INT NOT NULL,
    form_id INT NOT NULL,
    Accepted TINYINT NOT NULL,
    started_at DATE NOT NULL,

    UNIQUE KEY (user_id, event_id),

    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events(id) ON UPDATE CASCADE,
    FOREIGN KEY (form_id) REFERENCES forms(id) ON UPDATE CASCADE
);