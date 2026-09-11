DROP TABLE event_form;
DROP TABLE event_application;
DROP TAble event_participants;
DROP TABLE event_coords;
DROP TABLE events;

CREATE TABLE collaborators (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name_ar VARCHAR(255) NOT NULL,
    name_en VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    signature_id INT NOT NULL,

    FOREIGN KEY (signature_id) REFERENCES files(id) ON DELETE CASCADE
);

CREATE TABLE event (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    coordinator_id INT NOT NULL,
    presenter_id INT NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    max_participants INT NOT NULL,
    form_application TINYINT NOT NULL,
    outcomes TEXT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,

    FOREIGN KEY (coordinator_id) REFERENCES collaborators(id) ON UPDATE CASCADE,
    FOREIGN KEY (presenter_id) REFERENCES collaborators(id) ON UPDATE CASCADE
);

CREATE TABLE event_component (
    id INT PRIMARY KEY AUTO_INCREMENT,
    event_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    max_score DECIMAL(10, 2) NOT NULL,

    FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE
);

CREATE TABLE event_participant (
    id INT PRIMARY KEY AUTO_INCREMENT,
    event_id INT NOT NULL,
    user_id INT NOT NULL,
    grade DECIMAL(5, 2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    joined_at DATE NOT NULL,
    completed TINYINT NOT NULL,

    UNIQUE (event_id, user_id),

    FOREIGN KEY (event_id) REFERENCES event(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE component_score (
    id INT PRIMARY KEY AUTO_INCREMENT,
    participant_id INT NOT NULL,
    component_id INT NOT NULL,
    score DECIMAL(10, 2),

    UNIQUE (participant_id, component_id),

    FOREIGN KEY (participant_id) REFERENCES event_participant(id) ON DELETE CASCADE,
    FOREIGN KEY (component_id) REFERENCES event_component(id) ON DELETE CASCADE
);

CREATE TABLE certificate (
    id INT PRIMARY KEY AUTO_INCREMENT,
    cert_hash VARCHAR(255) UNIQUE NOT NULL,
    user_id INT NOT NULL,
    event_id INT NOT NULL,
    grade DECIMAL(5, 2) NOT NULL,
    issue_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) NOT NULL,

    UNIQUE (user_id, event_id),

    FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE,
    FOREIGN KEY (event_id) REFERENCES event(id)
);

CREATE TABLE certificate_file (
    id INT PRIMARY KEY AUTO_INCREMENT,
    certificate_id INT NOT NULL,
    store_id INT NOT NULL UNIQUE,
    lang VARCHAR(10) NOT NULL,

    UNIQUE (certificate_id, lang),

    FOREIGN KEY (certificate_id) REFERENCES certificate(id),
    FOREIGN KEY (store_id) REFERENCES files(id) ON DELETE CASCADE
);

CREATE TABLE event_form (
    id INT PRIMARY KEY AUTO_INCREMENT,
    form_id INT NOT NULL,
    event_id INT NOT NULL,
    
    UNIQUE KEY (form_id, event_id),

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